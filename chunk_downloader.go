package pdl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/fabienogli/pdl/chunker"
)

type simpleChunkDownloader struct {
	httpClient httpClient
}

func NewSimpleChunkDownloader(httpClient httpClient) *simpleChunkDownloader {
	return &simpleChunkDownloader{
		httpClient: httpClient,
	}
}

type chunkError struct {
	chunk chunker.Chunk
	err   error
}

func newChunkError(chunk chunker.Chunk, err error) chunkError {
	return chunkError{
		chunk: chunk,
		err:   err,
	}
}

func (p *simpleChunkDownloader) downloadPart(ctx context.Context, start, end int64, url, fileName string) error {
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

func (p *simpleChunkDownloader) wrappedDownload(ctx context.Context, id int, chunk chunker.Chunk, c chan chunkError, url, filename string) {
	err := p.downloadPart(ctx, chunk.Start, chunk.End, url, filename)
	if err != nil {
		err = fmt.Errorf("chunks[%d]: error while downloading part starting %d, ending at %d: %w", id, chunk.Start, chunk.End, err)
	}
	c <- newChunkError(chunk, err)
}

func (p *simpleChunkDownloader) Download(ctx context.Context, url, fileName string, chunks []chunker.Chunk) ([]chunker.Chunk, error) {
	var err error
	c := make(chan chunkError)
	for i, chunk := range chunks {
		go p.wrappedDownload(ctx, i, chunk, c, url, fileName)
	}
	var chunkLeft []chunker.Chunk
	for i := 0; i < len(chunks); i++ {
		chunkErr := <-c
		if chunkErr.err != nil {
			err = errors.Join(err, fmt.Errorf("chunk err: %w", chunkErr.err))
			chunkLeft = append(chunkLeft, chunkErr.chunk)
		}
	}
	return chunkLeft, err
}

type ChunkDownloaderUntilFailure struct {
	chunkDownloader chunkDownloader
	logger          logger
}

func NewChunkDownloaderUntilFailure(chunkDownloader chunkDownloader, logger logger) *ChunkDownloaderUntilFailure {
	return &ChunkDownloaderUntilFailure{
		chunkDownloader: chunkDownloader,
		logger:          logger,
	}
}

func (c *ChunkDownloaderUntilFailure) Download(ctx context.Context, url, fileName string, chunks []chunker.Chunk) ([]chunker.Chunk, error) {
	if ctx.Err() != nil {
		return chunks, fmt.Errorf("err ctx: %w", ctx.Err())
	}
	leftChunks, err := c.chunkDownloader.Download(ctx, url, fileName, chunks)
	if err != nil {
		c.logger.Printf("Retrying %d chunks for %s", len(leftChunks), fileName)
		return c.Download(ctx, url, fileName, leftChunks)
	}
	return nil, nil
}
