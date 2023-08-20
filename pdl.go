package pdl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/fabienogli/pdl/chunker"
)

type WrappedChunk struct {
	id    int
	chunk chunker.Chunk
	err   error
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type logger interface {
	Printf(format string, v ...any)
}

// ParallelDownloader parallel download
type ParallelDownloader struct {
	chunkSize  int64
	httpClient httpClient
	logger     logger
}

// NewParallelDownloader instatiate a ParallelDownloader
// chunkSize in Byte
func NewParallelDownloader(chunkSize int64, httpClient httpClient, logger logger) *ParallelDownloader {
	return &ParallelDownloader{
		chunkSize:  chunkSize,
		httpClient: httpClient,
		logger:     logger,
	}
}

func (p *ParallelDownloader) wrappedDownload(ctx context.Context, id int, chunk chunker.Chunk, c chan WrappedChunk, url, filename string) {
	err := p.downloadPart(ctx, chunk.Start, chunk.End, url, filename)
	if err != nil {
		err = fmt.Errorf("chunks[%d]: error while downloading part starting %d, ending at %d: %w", id, chunk.Start, chunk.End, err)
	}
	t := WrappedChunk{
		chunk: chunk,
		err:   err,
	}
	c <- t
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
	c := make(chan WrappedChunk)
	for i, chunk := range chunks {
		go p.wrappedDownload(ctx, i, chunk, c, url, fileName)
	}
	for i := 0; i < nbChunk; i++ {
		chunk := <-c
		if chunk.err != nil {
			err = errors.Join(err, fmt.Errorf("%w", chunk.err))
		}
	}
	return err
}

func (p *ParallelDownloader) downloadPart(ctx context.Context, start, end int64, url, fileName string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("err when http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := p.httpClient.Do(req)
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

func (p *ParallelDownloader) downloadSize(ctx context.Context, url string) (int64, error) {
	// GO idiom: acquire resource, check for error, defer release
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return 0, fmt.Errorf("err when http.NewRequestWithContext: %w", err)
	}
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
