package pdl

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fabienogli/pdl/chunker"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type logger interface {
	Printf(format string, v ...any)
}

type chunkDownloader interface {
	Download(ctx context.Context, url, fileName string, chunks []chunker.Chunk) ([]chunker.Chunk, error)
}

// ParallelDownloader parallel download
type ParallelDownloader struct {
	chunkSize       int64
	httpClient      httpClient
	logger          logger
	chunkDownloader chunkDownloader
}

// NewParallelDownloader instatiate a ParallelDownloader
// chunkSize in Byte
func NewParallelDownloader(chunkSize int64, httpClient httpClient, logger logger, chunkDownloader chunkDownloader) *ParallelDownloader {
	log.Println(httpClient)
	return &ParallelDownloader{
		chunkSize:       chunkSize,
		httpClient:      httpClient,
		logger:          logger,
		chunkDownloader: chunkDownloader,
	}
}

// Download downloads url into fileName
func (p *ParallelDownloader) Download(ctx context.Context, url, fileName string) error {
	size, err := p.downloadSize(ctx, url)
	if err != nil {
		p.logger.Printf("error happened downloading headers - %s", err)
		return err
	}
	p.logger.Printf("info url=%s: size=%d MB", url, size/1_000_000)
	chunks, err := chunker.Split(size, p.chunkSize)
	nbChunk := len(chunks)
	p.logger.Printf("nb of Chunk: %d", nbChunk)
	if err != nil {
		return err
	}
	err = createFile(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	_, err = p.chunkDownloader.Download(ctx, url, fileName, chunks)
	if err != nil {
		return fmt.Errorf("err when p.chunkDownloader.Download: %w", err)
	}
	return nil
}

func (p *ParallelDownloader) downloadSize(ctx context.Context, url string) (int64, error) {
	// GO idiom: acquire resource, check for error, defer release
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return 0, fmt.Errorf("err when http.NewRequestWithContext: %w", err)
	}
	log.Println(p.httpClient)
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("err when doing Head request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro status not ok: %s", resp.Status)
	}
	return resp.ContentLength, nil
}
