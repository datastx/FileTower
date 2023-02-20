package writers

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func StreamFileToS3(fileName, bucket, s3Location string, uploader *s3manager.Uploader) error {

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = io.Copy(pw, file)
		if err != nil {
			// put err in channel
			panic(err)
		}
	}()

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3Location),
		Body:   pr,
	})
	if err != nil {
		return err
	}
	return nil
}
