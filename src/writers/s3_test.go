package writers

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type mockUploader struct {
	uploadedData []byte
	err          error
}

func (u *mockUploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	if u.err != nil {
		return nil, u.err
	}

	// Read the data from the pipe and store it in the mock uploader
	data, err := ioutil.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	u.uploadedData = data

	return &s3manager.UploadOutput{}, nil
}

func TestStreamFileToS3(t *testing.T) {
	// Create a mock uploader
	mock := &mockUploader{}

	// Call the StreamFileToS3 function with the mock uploader
	err := StreamFileToS3("testdata/testfile.txt", "test-bucket", "testkey.txt", mock)

	if err != nil {
		t.Fatalf("StreamFileToS3 failed with error: %v", err)
	}

	// Verify that the data was uploaded correctly
	expectedData, err := ioutil.ReadFile("testdata/testfile.txt")
	if err != nil {
		t.Fatalf("Failed to read testdata file: %v", err)
	}

	if !bytes.Equal(mock.uploadedData, expectedData) {
		t.Errorf("Uploaded data does not match expected data")
	}
}
