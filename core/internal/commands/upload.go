package commands

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"

	"github.com/99xtal/bluthinator/core/internal/config"
)

type Frame struct {
	Episode   string
	Timestamp int
}

var uploadCmd = &cobra.Command{
	Use:   "upload [data_dir]",
	Short: "Upload extracted frames to object storage and the database",
	Long:  ``,
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
			Creds: credentials.NewStaticV4(config.ObjectStorageUser, config.ObjectStoragePass, ""),
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

		uploader := NewFrameUploader(db, minioClient, p)

		for i := uint(0); i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for episodePath := range episodeChan {
					err := uploader.UploadEpisode(episodePath)
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
