package controller

import (
	"context"

	"github.com/suzuki-shunsuke/ghir/pkg/github"
)

type Controller struct {
	input *Input
}

func New(input *Input) *Controller {
	return &Controller{
		input: input,
	}
}

type GitHub interface {
	ListReleases(ctx context.Context, owner, repo string) ([]*github.Release, error)
	EditRelease(ctx context.Context, owner, repo string, releaseID int64) error
}

type Input struct {
	GitHub GitHub
}

func NewInput() *Input {
	return &Input{}
}
