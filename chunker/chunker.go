package chunker

import "fmt"

type Chunk struct {
	Start int64
	End   int64
}

// Split return an array of Chunk of the len of chunkSize
func Split(size, chunkSize int64) ([]Chunk, error) {
	if size < 0 || chunkSize <= 0 {
		return nil, fmt.Errorf("size=%d, chunkSize=%d", size, chunkSize)
	}
	var cs []Chunk
	for start := int64(0); start < size; start += chunkSize {
		end := start + chunkSize - 1
		if end > size-1 {
			end = size - 1
		}
		c := Chunk{Start: start, End: end}
		cs = append(cs, c)
	}
	return cs, nil
}
