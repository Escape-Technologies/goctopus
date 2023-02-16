package utils

import (
	"bytes"
	"io"
	"os"
)

func CountLines(f *os.File) (int, error) {
	// input, err := os.Open(inputFile)
	// if err != nil {
	// 	log.Error(err)
	// 	os.Exit(1)
	// }
	// defer input.Close()
	
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
			c, err := f.Read(buf)
			count += bytes.Count(buf[:c], lineSep)

			switch {
			case err == io.EOF:
					// Seek is used to reset the file pointer to the beginning of the file
					f.Seek(0, io.SeekStart)
					return count, nil

			case err != nil:
					f.Seek(0, io.SeekStart)
					return count, err
			}
	}
}