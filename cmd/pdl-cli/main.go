package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fabienogli/pdl"
	"github.com/spf13/cobra"
)

// const url = "https://d37ci6vzurychx.cloudfront.net/trip-data/green_tripdata_2018-03.parquet"
// const outFile = "green_tripdata_2018-03.parquet"

// flags
const (
	timeout   = "timeout"
	chunkSize = "chunk-size"
)

func extractFilenameFromURL(url string) (string, error) {
	split := strings.Split(url, "/")
	if len(split) < 1 {
		return "", fmt.Errorf("url malformed: not '/' found in '%s'", url)
	}
	return split[len(split)-1], nil
}

func download(timeoutInSecond, chunkSizeInMB int64, url, outFile string, logger log.Logger) error {
	timeoutDuration := time.Duration(timeoutInSecond) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	httpClient := &http.Client{}
	chunkDownloader := pdl.NewSimpleChunkDownloader(httpClient)
	chunkDownloaderUntilFailure := pdl.NewChunkDownloaderUntilFailure(chunkDownloader, &logger)
	downloader := pdl.NewParallelDownloader(chunkSizeInMB*1_000_000, httpClient, &logger, chunkDownloaderUntilFailure)
	err := downloader.Download(ctx, url, outFile)
	if err != nil {
		return fmt.Errorf("err downloader.Download: %w", err)
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pdl-cli [url to download]",
	Short: "Pdl will download any file from a given URL",
	Long: `PDL (Parallel Downloader) is a CLI library that will download an URL.

It will chunk the file in several files
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		timeout, err := cmd.Flags().GetInt64(timeout)
		if err != nil {
			return fmt.Errorf("err getting timeout %w", err)
		}
		chunkSizeInMB, err := cmd.Flags().GetInt64(chunkSize)
		if err != nil {
			return fmt.Errorf("err getting chunk-size %w", err)
		}
		for _, url := range args {
			outFile, errFromLoop := extractFilenameFromURL(url)
			if errFromLoop != nil {
				err = errors.Join(err, fmt.Errorf("err extractFilenameFromURL: %w", errFromLoop))
				continue
			}
			errFromLoop = download(timeout, chunkSizeInMB, url, outFile, *log.Default())
			if errFromLoop != nil {
				err = errors.Join(err, fmt.Errorf("err download: %w", errFromLoop))
				continue
			}
			cmd.Printf("File %s was download\n", outFile)
		}
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().Int64P("timeout", "t", 60, "duration in seconds that the command will be canceled")
	rootCmd.PersistentFlags().Int64P("chunk-size", "s", 1, "the file will be split by this chunk size (MB)")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
