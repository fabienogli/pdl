package pdl

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/fabienogli/pdl/chunker"
	"github.com/fabienogli/pdl/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChunkDownloaderUntilFailure_Download(t *testing.T) {
	type fields struct {
		chunkDownloader chunkDownloader
	}
	type args struct {
		url      string
		fileName string
		chunks   []chunker.Chunk
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		cancelFunc bool
		wantErr    bool
	}{
		{
			name: "ok",
			fields: fields{
				chunkDownloader: func() chunkDownloader {
					m := mocks.NewChunkDownloader(t)
					return m
				}(),
			},
			args: args{
				url:      "https://superurl.com/file",
				fileName: "download_file",
				chunks: []chunker.Chunk{
					{Start: 0, End: 2},
					{Start: 3, End: 5},
					{Start: 6, End: 8},
					{Start: 9, End: 11},
				},
			},
			cancelFunc: true,
			wantErr:    true,
		},
		{
			name: "ok",
			fields: fields{
				chunkDownloader: func() chunkDownloader {
					m := mocks.NewChunkDownloader(t)
					m.On("Download", mock.Anything, "https://superurl.com/file", "download_file", []chunker.Chunk{
						{Start: 0, End: 2},
						{Start: 3, End: 5},
						{Start: 6, End: 8},
						{Start: 9, End: 11},
					}).Return([]chunker.Chunk{
						{Start: 0, End: 2},
						{Start: 6, End: 8},
					}, fmt.Errorf("error")).Once()
					m.On("Download", mock.Anything, "https://superurl.com/file", "download_file", []chunker.Chunk{
						{Start: 0, End: 2},
						{Start: 6, End: 8},
					}).Return([]chunker.Chunk{
						{Start: 6, End: 8},
					}, fmt.Errorf("error")).Once()
					m.On("Download", mock.Anything, "https://superurl.com/file", "download_file", []chunker.Chunk{
						{Start: 6, End: 8},
					}).Return(nil, nil).Once()

					return m
				}(),
			},
			args: args{
				url:      "https://superurl.com/file",
				fileName: "download_file",
				chunks: []chunker.Chunk{
					{Start: 0, End: 2},
					{Start: 3, End: 5},
					{Start: 6, End: 8},
					{Start: 9, End: 11},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancelFunc()
			if tt.cancelFunc {
				cancelFunc()
			}
			c := &ChunkDownloaderUntilFailure{
				chunkDownloader: tt.fields.chunkDownloader,
				logger:          log.Default(),
			}
			_, err := c.Download(ctx, tt.args.url, tt.args.fileName, tt.args.chunks)
			assert.Equal(t, tt.wantErr, err != nil, fmt.Sprintf("ChunkDownloaderUntilFailure.Download() error = %v, wantErr %v", err, tt.wantErr))
		})
	}
}
