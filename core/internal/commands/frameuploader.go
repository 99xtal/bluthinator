package commands

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/vbauerster/mpb/v7"
)

type FrameUploader struct {
	db          *sql.DB
	minioClient *minio.Client
	p           *mpb.Progress
}

func NewFrameUploader(db *sql.DB, minioClient *minio.Client, p *mpb.Progress) *FrameUploader {
	return &FrameUploader{
		db:          db,
		minioClient: minioClient,
		p:           p,
	}
}

func (fu *FrameUploader) UploadEpisode(episodeDir string) error {
	episode := filepath.Base(episodeDir)
	frameDir := filepath.Join(episodeDir, "frames")

	updatedCount, err := fu.syncFilesWithStorage(episode, frameDir)
	if err != nil {
		return err
	}
	if updatedCount == 0 {
		fmt.Printf("[%s] Files already synced in object storage\n", episode)
		return nil
	}

	fmt.Printf("[%s] Reading frame index\n", episode)
	frames, err := fu.readFrameIndexCSV(frameDir)
	if err != nil {
		return err
	}

	err = fu.rebuildDBIndex(episode, frames)
	if err != nil {
		return err
	}

	return nil
}

func (fu *FrameUploader) syncFilesWithStorage(episodeKey string, frameDir string) (int, error) {
	bucketName := "bluthinator"
	objPrefix := fmt.Sprintf("frames/%s", episodeKey)
	ctx := context.Background()

	fmt.Printf("[%s] Comparing local and remote files\n", episodeKey)

	var files []string
	err := filepath.Walk(frameDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	addedCount := 0
	removedCount := 0

	// List objects in Minio
	objectCh := fu.minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    objPrefix,
		Recursive: true,
	})

	minioObjects := make(map[string]struct{})
	for object := range objectCh {
		if object.Err != nil {
			return 0, object.Err
		}
		minioObjects[object.Key] = struct{}{}
	}

	// Pre-process to determine files to add and delete
	localFiles := make(map[string]struct{})
	for _, filePath := range files {
		relativePath, err := filepath.Rel(frameDir, filePath)
		if err != nil {
			return 0, err
		}

		objectKey := fmt.Sprintf("frames/%s/%s", episodeKey, relativePath)
		localFiles[objectKey] = struct{}{}

		if _, exists := minioObjects[objectKey]; !exists {
			addedCount++
		}
	}
	updateTotal := addedCount + removedCount

	if updateTotal == 0 {
		return 0, nil
	}

	bar := newProgressBar(fu.p, int64(updateTotal), fmt.Sprintf("[%s] Syncing local and remote files in object storage ", episodeKey))

	for _, filePath := range files {
		start := time.Now()

		relativePath, err := filepath.Rel(frameDir, filePath)
		if err != nil {
			return 0, err
		}

		objectKey := fmt.Sprintf("frames/%s/%s", episodeKey, relativePath)
		if _, exists := minioObjects[objectKey]; !exists {
			addedCount++
			_, err := fu.minioClient.FPutObject(ctx, bucketName, objectKey, filePath, minio.PutObjectOptions{})
			if err != nil {
				return 0, err
			}
			bar.Increment()
			bar.DecoratorEwmaUpdate(time.Since(start))
		}
	}

	// Delete objects from Minio that aren't present locally
	for objectName := range minioObjects {
		if _, exists := localFiles[objectName]; !exists {
			err := fu.minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
			if err != nil {
				return 0, err
			}
			bar.Increment()
		}
	}

	bar.SetTotal(bar.Current(), true)
	bar.Wait()
	fmt.Printf("[%s] Synced %d files (%d added, %d removed)\n", episodeKey, len(files), addedCount, removedCount)

	return updateTotal, nil
}

func (fu *FrameUploader) readFrameIndexCSV(frameDir string) ([]Frame, error) {
	csvFilePath := filepath.Join(frameDir, "index.csv")

	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	frames := make([]Frame, 0, len(records))
	for _, record := range records[1:] {
		if len(record) != 2 {
			continue
		}

		timestamp, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}

		frame := Frame{
			Episode:   record[0],
			Timestamp: timestamp,
		}
		frames = append(frames, frame)
	}

	return frames, nil
}

func (fu *FrameUploader) rebuildDBIndex(episode string, frames []Frame) error {
	fmt.Printf("[%s] Rebuilding frame index in DB\n", episode)
	// Delete existing frames from DB
	result, err := fu.db.Exec("DELETE FROM frames WHERE episode=$1", episode)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		fmt.Printf("[%s] Deleted %d rows from frames table\n", episode, rowsAffected)
	}

	// Upload frame index to DB using batch insert
	bar := newProgressBar(fu.p, int64(len(frames)), fmt.Sprintf("[%s] Uploading frame index to DB ", episode))

	const batchSize = 1000
	for i := 0; i < len(frames); i += batchSize {
		start := time.Now()

		end := i + batchSize
		if end > len(frames) {
			end = len(frames)
		}

		batch := frames[i:end]
		values := make([]interface{}, 0, len(batch)*2)
		placeholders := make([]string, 0, len(batch))

		for j, frame := range batch {
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", j*2+1, j*2+2))
			values = append(values, frame.Episode, frame.Timestamp)
		}

		query := fmt.Sprintf("INSERT INTO frames (episode, timestamp) VALUES %s", strings.Join(placeholders, ","))
		_, err := fu.db.Exec(query, values...)
		if err != nil {
			return err
		}

		bar.SetCurrent(int64(end))
		bar.DecoratorEwmaUpdate(time.Since(start))
	}

	bar.SetTotal(bar.Current(), true)
	bar.Wait()

	return nil
}
