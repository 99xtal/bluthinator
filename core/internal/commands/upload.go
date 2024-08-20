package commands

import (
	"context"
	"database/sql"
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

func uploadEpisodeFrames(episodeFramesPath string, p *mpb.Progress, db *sql.DB, minioClient *minio.Client) error {
	episode := filepath.Base(episodeFramesPath)

	updatedCount, err := syncFilesWithStorage(episodeFramesPath, p, minioClient)
	if err != nil {
		return err
	}
	if updatedCount == 0 {
		fmt.Printf("[%s] Files already synced in object storage\n", episode)
		return nil
	}

	err = rebuildDBIndex(episodeFramesPath, p, db)
	if err != nil {
		return err
	}

	return nil
}

func syncFilesWithStorage(episodeFramesPath string, p *mpb.Progress, minioClient *minio.Client) (int, error) {
	episodeKey := filepath.Base(episodeFramesPath)
	bucketName := "bluthinator"
	prefix := fmt.Sprintf("frames/%s", episodeKey)
	ctx := context.Background()

	var files []string
	err := filepath.Walk(episodeFramesPath, func(path string, info os.FileInfo, err error) error {
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
		Prefix:    prefix,
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
		relativePath, err := filepath.Rel(episodeFramesPath, filePath)
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

	bar := newProgressBar(p, int64(addedCount+removedCount), fmt.Sprintf("[%s] Syncing local files with files in object storage ", episodeKey))

	for _, filePath := range files {
		start := time.Now()

		relativePath, err := filepath.Rel(episodeFramesPath, filePath)
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

	bar.Wait()
	fmt.Printf("[%s] Synced %d files (%d added, %d removed)\n", episodeKey, len(files), addedCount, removedCount)

	return updateTotal, nil
}

func rebuildDBIndex(episodeFramesPath string, p *mpb.Progress, db *sql.DB) error {
	episodeKey := filepath.Base(episodeFramesPath)

	// Delete existing frames from DB
	result, err := db.Exec("DELETE FROM frames WHERE episode=$1", episodeKey)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		fmt.Printf("[%s] Deleted %d rows from frames table\n", episodeKey, rowsAffected)
	}

	// Create index of frames
	var frames []Frame;
	err = filepath.Walk(episodeFramesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == episodeFramesPath {
			return nil
		}

		if info.IsDir() {
			relativePath, err := filepath.Rel(episodeFramesPath, path)
			if err != nil {
				return err
			}

			pathParts := strings.Split(relativePath, string(os.PathSeparator))
			if len(pathParts) == 1 {
				timestamp, err := strconv.ParseInt(pathParts[0], 10, 64)
				if err != nil {
					return err
				}

				frame := Frame{
					Episode: episodeKey,
					Timestamp: int(timestamp),
				}
				frames = append(frames, frame)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Upload frame index to DB using batch insert
	bar := newProgressBar(p, int64(len(frames)), fmt.Sprintf("[%s] Uploading frame index to DB ", episodeKey))

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