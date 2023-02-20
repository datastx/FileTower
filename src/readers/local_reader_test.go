package readers

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadInChunks(t *testing.T) {
	// Create a test file
	filename := "test.txt"
	data := []byte("This is a test file.")
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Read the test file in chunks and verify the data
	var result []byte
	processData := func(data []byte) {
		result = append(result, data...)
	}
	err = readInChunks(filename, processData)
	if err != nil {
		t.Fatalf("Failed to read file in chunks: %v", err)
	}
	if !bytes.Equal(result, data) {
		t.Errorf("Expected result to be %q but got %q", data, result)
	}

	// Clean up the test file
	err = os.Remove(filename)
	if err != nil {
		t.Fatalf("Failed to clean up test file: %v", err)
	}
}
