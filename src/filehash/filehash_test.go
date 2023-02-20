package filehash

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetCheckSum(t *testing.T) {
	// Create a temporary file
	file, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write some data to the file
	_, err = file.WriteString("hello world")
	if err != nil {
		t.Fatal(err)
	}

	// Calculate the expected checksum
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9" // SHA256("hello world")

	// Calculate the actual checksum using the function being tested
	actual := getCheckSum(file.Name())

	// Compare the expected and actual checksums
	if actual != expected {
		t.Errorf("getCheckSum(%q) = %q, expected %q", file.Name(), actual, expected)
	}
}
