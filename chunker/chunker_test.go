package chunker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	type args struct {
		size      int64
		chunkSize int64
	}
	tests := []struct {
		name    string
		args    args
		want    []Chunk
		wantErr bool
	}{
		{
			name: "nil result",
			args: args{
				size:      0,
				chunkSize: 1,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "working case",
			args: args{
				size:      13,
				chunkSize: 3,
			},
			want: []Chunk{
				{0, 2}, {3, 5}, {6, 8}, {9, 11}, {12, 12},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Split(tt.args.size, tt.args.chunkSize)
			require.Equal(t, tt.wantErr, err != nil, fmt.Sprintf("want err %t, got %v", tt.wantErr, err))
			assert.Equal(t, tt.want, got)
		})
	}
}
