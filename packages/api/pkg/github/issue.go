package github

import (
	"context"

	"github.com/google/go-github/v53/github"
	"github.com/pkg/errors"
)

func (git *_github) CreateIssue(ctx context.Context, issue github.IssueRequest) (*int, error) {
	if ctx == nil || issue == (github.IssueRequest{}) {
		return nil, errors.New("Input is empty")
	}
	i, _, err := git.client.Issues.Create(ctx, git.owner, git.repo, &issue)
	if err != nil {
		return nil, err
	}
	return i.Number, nil
}

func (git *_github) UpdateIssue(ctx context.Context, number int, issue github.IssueRequest) error {
	if ctx == nil || issue == (github.IssueRequest{}) {
		return errors.New("Input is empty")
	}
	_, _, err := git.client.Issues.Edit(ctx, git.owner, git.repo, number, &issue)
	if err != nil {
		return err
	}
	return nil
}
