package storage

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsS3 struct {
	awss3       *s3.S3
	bucketRead  string
	bucketWrite string
	bucketSQS   string
	urlQueue    string
}

func NewS3Storage() (*AwsS3, error) {

	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	awss3 := s3.New(session)

	conn := &AwsS3{
		awss3:       awss3,
		bucketRead:  os.Getenv("BUCKET_READ_NAME"),
		bucketWrite: os.Getenv("BUCKET_WRITE_NAME"),
		bucketSQS:   *aws.String("jokes-sqs-messages"),
		urlQueue:    *aws.String("https://sqs.eu-central-1.amazonaws.com/333746971525/JokesQueueSend"),
	}
	return conn, nil
}

func (a *AwsS3) UploadMessageTos3(name string, data io.Reader) error {

	uploader := s3manager.NewUploaderWithClient(a.awss3)
	upParams := &s3manager.UploadInput{
		Bucket: &a.bucketSQS,
		Key:    &name,
		Body:   data,
	}
	_, err := uploader.Upload(upParams, func(u *s3manager.Uploader) {
		u.LeavePartsOnError = true
	})
	if err != nil {
		return err
	}

	return nil
}
