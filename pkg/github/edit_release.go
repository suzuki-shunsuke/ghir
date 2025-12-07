package github

import (
	"context"

	"github.com/google/go-github/v80/github"
)

func (c *Client) EditRelease(ctx context.Context, owner, repo string, id int64) error {
	// GraphQL API does not support updating a release.
	// https://docs.github.com/en/graphql/reference/mutations
	// https://pkg.go.dev/github.com/google/go-github/v75/github#RepositoriesService.EditRelease
	_, _, err := c.repos.EditRelease(ctx, owner, repo, id, &github.RepositoryRelease{})
	return err //nolint:wrapcheck
}
