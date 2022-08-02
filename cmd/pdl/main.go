package main

import (
	"context"
	"log"
	"time"

	"github.com/fabienogli/pdl"
)

const url = "https://d37ci6vzurychx.cloudfront.net/trip-data/green_tripdata_2018-03.parquet"
const outFile = "green_tripdata_2018-03.parquet"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := pdl.Download(ctx, url, outFile)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
