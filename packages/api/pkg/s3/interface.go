package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Service struct {
	s3 *s3.S3
}

type S3Service interface {
	DownloadFile(bucket, key string) (io.ReadCloser, error)
}
