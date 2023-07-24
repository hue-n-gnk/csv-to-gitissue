package s3

import (
	"fmt"
	"path"
	"time"
)

func CreateS3Key() string {
	now := time.Now().Format("20060102")
	return path.Join(fmt.Sprintf("Backlog-Issues-%s.csv", now))
}
