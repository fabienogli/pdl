package pdl

import (
	"bytes"
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
	log.Printf("info url=%s: size=%d MB", url, size/1_000_000)
	chunks, err := chunker.Split(size, 1_000_000)
	log.Printf("nb of Chunk: %d", len(chunks))
	if err != nil {
		return err
	}
	err = createFile(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	var wg sync.WaitGroup
	wg.Add(len(chunks))
	for i, c := range chunks {
		go func(id int, chunk chunker.Chunk) {
			defer wg.Done()
			log.Printf("chunks[%d] started", id)
			if err := downloadPart(ctx, chunk.Start, chunk.End, url, fileName); err != nil {
				log.Printf("chunks[%d]: error while downloading part starting %d, ending at %d: %v", id, chunk.Start, chunk.End, err)
				return
			}
			log.Printf("chunks[%d] success", id)
		}(i, c)
	}
	wg.Wait()
	return nil
}

func downloadPart(ctx context.Context, start, end int64, url, fileName string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("err when http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("err when downloading: %w", err)
	}
	defer resp.Body.Close()
	buf := bytes.NewBuffer([]byte{})
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("err when buf.ReadFrom: %w", err)
	}
	err = writeAt(fileName, start, buf.Bytes())
	if err != nil {
		return fmt.Errorf("err when writeAt: %w", err)
	}
	return nil
}

func downloadSize(ctx context.Context, url string) (int64, error) {
	// GO idiom: acquire resource, check for error, defer release
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return 0, fmt.Errorf("err when http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("err when doing Head request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro status not ok: %s", resp.Status)
	}
	return resp.ContentLength, nil
}
