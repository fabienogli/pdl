package pdl

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/fabienogli/pdl/chunker"
)

// Download downloads url into fileName
func Download(ctx context.Context, url, fileName string) error {
	size, err := downloadSize(ctx, url)
	if err != nil {
		log.Printf("error happened downloading headers - %s", err)
		return err
	}
	log.Printf("info %#v: size=%d", url, size)
	chunks, err := chunker.Split(size, 1_000_000)
	if err != nil {
		return err
	}
	err = createFile(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	// os.Create("")
	var wg sync.WaitGroup
	wg.Add(len(chunks))
	for _, c := range chunks {
		go func(chunk chunker.Chunk) {
			defer wg.Done()
			if err := downloadPart(ctx, chunk.Start, chunk.End, url, fileName); err != nil {
				log.Printf("error while downloading part starting %d, ending at %d", chunk.Start, chunk.End)
			}
		}(c)
	}
	return nil
}

func downloadPart(ctx context.Context, start, end int64, url, fileName string) error {
	//download date from http
	var data []byte
	return writeAt(fileName, start, data[start:end])
}

func downloadSize(ctx context.Context, url string) (int64, error) {
	// GO idiom: acquire resource, check for error, defer release
	req, err := http.NewRequestWithContext(ctx, url, "HEAD", nil)
	if err != nil {
		return 0, err
	}
	// resp, err := http.Head(url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	log.Println("Status code is ", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(resp.Status)
	}
	// resp.Header.Get("Content-Type") --> Header type is case insensitive
	return resp.ContentLength, nil
}
