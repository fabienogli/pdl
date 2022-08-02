package pdl

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func goodServer(expected int64) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(int(expected)))
	}))
}
func Test_downloadSize(t *testing.T) {
	const expected int64 = 500
	server := goodServer(expected)
	// http call should be mock
	size, err := downloadSize(context.Background(), server.URL)
	require.NoErrorf(t, err, server.URL)
	require.Equal(t, expected, size)
}

func Fuzz_downloadSize(f *testing.F) {
	log.Println("fuzzing")

	f.Fuzz(func(t *testing.T, expected int64) {
		server := goodServer(expected)
		defer server.Close()
		size, err := downloadSize(context.Background(), server.URL)
		require.NoErrorf(t, err, server.URL)
		require.Equal(t, expected, size)
	})
}
