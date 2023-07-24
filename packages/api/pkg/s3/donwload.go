package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *s3Service) DownloadFile(bucket, key string) (io.ReadCloser, error) {
	payload := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	object, err := s.s3.GetObject(payload)
	if err != nil {
		return nil, err
	}
	return object.Body, nil

}
