package s3

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewS3Service(sess client.ConfigProvider) S3Service {
	return &s3Service{
		s3: s3.New(sess),
	}
}
