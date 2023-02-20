package readers

import (
	"fmt"
	"io"
	"os"
)

func readInChunks(filename string, processData func(data []byte)) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a buffer to hold each chunk of data
	buf := make([]byte, 1024)

	for {
		// Read a chunk of data from the file into the buffer
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// Pass the data in the buffer to the processData callback
		processData(buf[:n])
	}

	return nil
}

func processData(data []byte) {
	// This function processes each chunk of data
	fmt.Printf("Read %d bytes: %s\n", len(data), string(data))
}
