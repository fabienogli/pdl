package chunker

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/gogo/protobuf/plugin/size"
// )

// func TestSplit(t *testing.T) {
// 	splitTestCases = []struct{
// 		size	int64
// 		chunkSize	int64
// 		chunks []Chunk
// 	} {
// 		{0, 1, nil},
// 		{6, 3, []Chunk]
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := Split(tt.args.size, tt.args.chunkSize)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Split() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
