package commands

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"

	"github.com/99xtal/bluthinator/core/internal/config"
)

type Frame struct {
	Episode string
	Timestamp int
}

var uploadCmd = &cobra.Command{
	Use:   "upload [frame_dir]",
	Short: "Upload extracted frames to object storage and the database",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		frameDir := args[0]
		config := config.New()

		connString := config.GetPostgresConnString()
		db, err := sql.Open("postgres", connString)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		minioClient, err := minio.New(config.ObjectStorageEndpoint, &minio.Options{
			Creds: credentials.NewStaticV4("minio", "minio123", ""),
		})
		if err != nil {
			log.Fatal(err)
		}

		episodeDirs, err := filepath.Glob(filepath.Join(frameDir, "*"))
		if err != nil {
			log.Fatal(err)
		}
		if len(episodeDirs) == 0 {
			fmt.Println("No episode directories found in the input directory")
			os.Exit(1)
		}

		episodeChan := make(chan string, len(episodeDirs))
		var wg sync.WaitGroup

		p := mpb.New(mpb.WithWaitGroup(&wg))

		for i := uint(0); i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for episodePath := range episodeChan {
					err := uploadEpisodeFrames(episodePath, p, db, minioClient)
					if err != nil {
						log.Fatal(err)
					}
				}
			}()
		}

		for _, episodePath := range episodeDirs {
			episodeChan <- episodePath
		}
		close(episodeChan)

		p.Wait()
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().UintVarP(&numWorkers, "workers", "w", 3, "Number of workers to use")
}

func uploadEpisodeFrames(episodeDir string, p *mpb.Progress, db *sql.DB, minioClient *minio.Client) error {
	episode := filepath.Base(episodeDir)
	frameDir := filepath.Join(episodeDir, "frames")

	updatedCount, err := syncFilesWithStorage(episodeDir, p, minioClient)
	if err != nil {
		return err
	}
	if updatedCount == 0 {
		fmt.Printf("[%s] Files already synced in object storage\n", episode)
		return nil
	}

	fmt.Printf("[%s] Reading frame index\n", episode)
	frames, err := readFrameIndexCSV(frameDir)
	if err != nil {
		return err
	}

	err = rebuildDBIndex(episode, frames, p, db)
	if err != nil {
		return err
	}

	return nil
}

func syncFilesWithStorage(episodeDir string, p *mpb.Progress, minioClient *minio.Client) (int, error) {
	episodeKey := filepath.Base(episodeDir)
	bucketName := "bluthinator"
	objPrefix := fmt.Sprintf("frames/%s", episodeKey)
	ctx := context.Background()

	frameDir := filepath.Join(episodeDir, "frames")

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
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
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

	bar := newProgressBar(p, int64(updateTotal), fmt.Sprintf("[%s] Syncing local and remote files in object storage ", episodeKey))

	for _, filePath := range files {
		start := time.Now()

		relativePath, err := filepath.Rel(frameDir, filePath)
        if err != nil {
            return 0, err
        }

		objectKey := fmt.Sprintf("frames/%s/%s", episodeKey, relativePath)
		if _, exists := minioObjects[objectKey]; !exists {
			addedCount++
			_, err := minioClient.FPutObject(ctx, bucketName, objectKey, filePath, minio.PutObjectOptions{})
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
			err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
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

func rebuildDBIndex(episode string, frames []Frame, p *mpb.Progress, db *sql.DB) error {
	fmt.Printf("[%s] Rebuilding frame index in DB\n", episode)
	// Delete existing frames from DB
	result, err := db.Exec("DELETE FROM frames WHERE episode=$1", episode)
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
	bar := newProgressBar(p, int64(len(frames)), fmt.Sprintf("[%s] Uploading frame index to DB ", episode))

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
		_, err := db.Exec(query, values...)
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

func readFrameIndexCSV(frameDir string) ([]Frame, error) {
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
			Episode: record[0],
			Timestamp: timestamp,
		}
		frames = append(frames, frame)
	}

	return frames, nil
}