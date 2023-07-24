package backlogcvs

import (
	"encoding/csv"
	"hue-n-gnk/csv-to-gitisue/helpers/env"
	"hue-n-gnk/csv-to-gitisue/pkg/s3"
	"io"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

func download() ([]CsvStruct, error) {
	e := env.Load()
	s3key := s3.CreateS3Key()
	sess, err := env.GetSession()
	if err != nil {
		return nil, err
	}
	s3Svc := s3.NewS3Service(sess)
	data, err := s3Svc.DownloadFile(e.BucketId, s3key)
	if err != nil {
		return nil, errors.Wrap(err, "LogsService: download file failed")
	}

	return csvToStruct(data)
}

func csvToStruct(reader io.ReadCloser) (data []CsvStruct, err error) {
	defer reader.Close()
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	err = gocsv.UnmarshalCSV(csvReader, &data)
	return
}
