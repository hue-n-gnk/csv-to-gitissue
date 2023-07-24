package backlogcvs

import (
	"hue-n-gnk/csv-to-gitisue/database/dao"
	"hue-n-gnk/csv-to-gitisue/database/factory"
	"hue-n-gnk/csv-to-gitisue/pkg/logger"
)

type BacklogCsvService interface {
	ToGithub() error
}

type CsvStruct struct {
	Id          string `csv:"ID"`
	KeyId       string `csv:"Key ID"`
	Key         string `csv:"Key"`
	IssueType   string `csv:"Issue Type"`
	Subject     string `csv:"Subject"`
	Description string `csv:"Description"`
	StatusId    string `csv:"Status ID"`
	Status      string `csv:"Status"`
	PriorityId  string `csv:"Priority ID"`
	Priority    string `csv:"Priority"`
	AssigneeID  string `csv:"Assignee ID"`
}

type backlogCsvService struct {
	db  *dao.Query
	log logger.Logger
}

func NewService() BacklogCsvService {
	db, err := factory.Query()
	if err != nil {
		panic(err)
	}
	return &backlogCsvService{
		db:  db,
		log: logger.NewAPILogger(),
	}
}
