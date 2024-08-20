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

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().UintVarP(&numWorkers, "workers", "w", 3, "Number of workers to use")
}

func uploadEpisodeFrames(episodeFramesPath string, p *mpb.Progress, db *sql.DB, minioClient *minio.Client) error {
	episode := filepath.Base(episodeFramesPath)

	result, err := db.Exec("DELETE FROM frames WHERE episode=$1", episode)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected > 0 {
		fmt.Printf("[%s] Deleted %d rows from frames table\n", episode, rowsAffected)
	}

	// List existing objects in MinIO
	prefix := fmt.Sprintf("frames/%s", episode)
	objectCh := minioClient.ListObjects(context.Background(), "bluthinator", minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var objects []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		objects = append(objects, object)
	}

	if len(objects) > 0 {
		bar := newProgressBar(p, int64(len(objects)), fmt.Sprintf("[%s] Deleting files from %s ", episode, prefix))

		for _, object := range objects {
			start := time.Now()

			err := minioClient.RemoveObject(context.Background(), "bluthinator", object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				log.Fatalln(err)
			}

			bar.Increment()
			bar.DecoratorEwmaUpdate(time.Since(start))
		}

		bar.SetTotal(bar.Current(), true)
		bar.Wait()
	}

	// Collect all files to be uploaded
	var files []string
	err = filepath.Walk(episodeFramesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Upload files to MinIO
	bar := newProgressBar(p, int64(len(files)), fmt.Sprintf("[%s] Uploading files to %s ", episode, prefix))

	var frames []Frame;
	for _, filePath := range files {
		start := time.Now()

		relativePath, err := filepath.Rel(episodeFramesPath, filePath)
        if err != nil {
            log.Fatalf("Failed to compute relative path: %v", err)
        }

        // Split the relative path to extract episode and timestamp
        pathParts := strings.Split(filePath, string(os.PathSeparator))
        if len(pathParts) < 2 {
            log.Fatalf("Invalid file path structure: %s", filePath)
        }
        episode := pathParts[1]
        timestamp, err := strconv.ParseInt(pathParts[2], 10, 64)
		if err != nil {
			return err
		}

        objectName := fmt.Sprintf("frames/%s/%s", episode, relativePath)

		// Upload file
		_, err = minioClient.FPutObject(context.Background(), "bluthinator", objectName, filePath, minio.PutObjectOptions{})
		if err != nil {
			return err
		}

		frame := Frame{
			Episode: episode,
			Timestamp: int(timestamp),
		}

		frames = append(frames, frame)

		bar.Increment()
		bar.DecoratorEwmaUpdate(time.Since(start))
	}

	bar.SetTotal(bar.Current(), true)
	bar.Wait()

	// Upload frame index to DB using batch insert
	bar = newProgressBar(p, int64(len(frames)), fmt.Sprintf("[%s] Uploading frame index to DB ", episode))

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