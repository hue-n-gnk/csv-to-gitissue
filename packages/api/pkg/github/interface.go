package github

import (
	"context"
	"hue-n-gnk/csv-to-gitisue/pkg/logger"

	"github.com/google/go-github/v53/github"
)

type GitHub interface {
	CreateIssue(ctx context.Context, issue github.IssueRequest) (*int, error)
	UpdateIssue(ctx context.Context, number int, issue github.IssueRequest) error
}

type _github struct {
	client *github.Client
	log    logger.Logger
	owner  string
	repo   string
}
