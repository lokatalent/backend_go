package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/lokatalent/backend_go/internal/repository"
)

type storage struct {
	bucket  string
	session *session.Session
	client  *s3.Client
}

func NewStorageInfrastructure(bucket string) repository.StorageRepository {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	client := s3.NewFromConfig(cfg)

	return storage{session: sess, bucket: bucket, client: client}
}

func (s storage) UploadFile(file io.Reader, key string, contentType string) (string, error) {
	uploader := s3manager.NewUploader(s.session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:        file,
		Bucket:      &s.bucket,
		Key:         &key,
		ContentType: &contentType,
	})

	if err != nil {
		return "", err
	}
	return result.Location, err
}

func (s storage) DeleteFile(key string) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}

	return nil
}
