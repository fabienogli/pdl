package pdl

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getRangeData(rangeData string, data []byte) ([]byte, error) {
	split := strings.Split(rangeData, "=")
	if len(split) < 2 {
		return nil, fmt.Errorf("err when strings.Split(%s)", rangeData)
	}
	newSplit := strings.Split(split[1], "-")
	if len(newSplit) < 2 {
		return nil, fmt.Errorf("err when strings.Split(%s)", split[1])
	}
	start, err := strconv.Atoi(newSplit[0])
	if err != nil {
		return nil, fmt.Errorf("err atoi start %s", newSplit[0])
	}
	end, err := strconv.Atoi(newSplit[1])
	if err != nil {
		return nil, fmt.Errorf("err atoi end %s", newSplit[1])
	}
	size := len(data)
	if start > size {
		return nil, fmt.Errorf("err start > size (%d > %d)", start, size)
	}
	if end > size {
		return nil, fmt.Errorf("err end > size (%d > %d)", end, size)
	}
	return data[start:end], nil
}

func goodServer(expected int64, data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Set("Content-Length", strconv.Itoa(int(expected)))
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			body, err := getRangeData(r.Header.Get("Range"), data)
			if err != nil {
				log.Printf("err getRangeData: %v", err)
				return
			}
			_, err = w.Write(body)
			if err != nil {
				log.Printf("err %v", err)
			}
			return
		}
	}))
}

func getExpected(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("err os.Open: %w", err)
	}
	return io.ReadAll(f)
}

func Test_Downlaod(t *testing.T) {

	expected, err := getExpected("expected.parquet")
	require.NoError(t, err)

	dir := t.TempDir()

	filename := fmt.Sprintf("%s/download.parquet", dir)

	server := goodServer(int64(len(expected)), expected)

	httpClient := &http.Client{}

	chunkDownLoader := NewSimpleChunkDownloader(httpClient)

	downloader := NewParallelDownloader(1_000_000, httpClient, log.Default(), chunkDownLoader)

	err = downloader.Download(context.Background(), server.URL, filename)
	require.NoError(t, err)

	got, err := os.ReadFile(filename)
	require.NoError(t, err)

	_ = got
	// assert.Equal(t, expected, got)

}

func Test_downloadSize(t *testing.T) {
	const expected int64 = 500
	server := goodServer(expected, nil)
	downloader := NewParallelDownloader(1_000_000, http.DefaultClient, log.Default(), nil)

	size, err := downloader.downloadSize(context.Background(), server.URL)
	require.NoErrorf(t, err, server.URL)
	require.Equal(t, expected, size)
}

func Fuzz_downloadSize(f *testing.F) {
	f.Fuzz(func(t *testing.T, expected int64) {
		if expected < 0 {
			expected = -expected
		}
		server := goodServer(expected, nil)
		defer server.Close()

		httpClient := http.DefaultClient

		chunkDownLoader := NewSimpleChunkDownloader(httpClient)

		downloader := NewParallelDownloader(1_000_000, httpClient, log.Default(), chunkDownLoader)

		size, err := downloader.downloadSize(context.Background(), server.URL)
		require.NoErrorf(t, err, server.URL)
		require.Equal(t, expected, size)
	})
}
