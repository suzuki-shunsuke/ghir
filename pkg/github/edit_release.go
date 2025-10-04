package github

import (
	"context"

	"github.com/google/go-github/v75/github"
)

type InputEditRelease struct {
	ID int64
}

func (c *Client) EditRelease(ctx context.Context, owner, repo string, id int64, description string) error {
	// GraphQL API does not support updating a release.
	// https://docs.github.com/en/graphql/reference/mutations
	// https://pkg.go.dev/github.com/google/go-github/v75/github#RepositoriesService.EditRelease
	_, _, err := c.repos.EditRelease(ctx, owner, repo, id, &github.RepositoryRelease{
		Body: github.Ptr(description),
	})
	return err //nolint:wrapcheck
}
