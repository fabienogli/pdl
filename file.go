package pdl

import (
	"fmt"
	"log"
	"os"
)

func createFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("err os.Create: %w", err)
	}
	defer f.Close()
	return nil
}

func writeAt(fileName string, offset int64, data []byte) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("err when os.OpenFile: %w", err)
	}
	defer f.Close()
	log.Printf("len %d", len(data))
	n, err := f.WriteAt(data, offset)
	log.Printf("wrote %d", n)
	if err != nil {
		return fmt.Errorf("err when f.WriteAt: %w", err)
	}
	return nil
}
