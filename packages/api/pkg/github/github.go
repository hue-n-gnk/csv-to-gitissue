package github

import (
	"context"
	"hue-n-gnk/csv-to-gitisue/helpers/env"
	"hue-n-gnk/csv-to-gitisue/pkg/logger"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func New() GitHub {
	e := env.Load()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: e.GHToken},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &_github{
		client: client,
		log:    logger.NewAPILogger(),
		owner:  e.GHOwner,
		repo:   e.GHRepo,
	}
}
