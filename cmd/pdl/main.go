package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/fabienogli/pdl"
)

const url = "https://d37ci6vzurychx.cloudfront.net/trip-data/green_tripdata_2018-03.parquet"
const outFile = "green_tripdata_2018-03.parquet"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	downloader := pdl.NewParallelDownloader(1_000_000, http.DefaultClient, log.Default())
	err := downloader.Download(ctx, url, outFile)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
