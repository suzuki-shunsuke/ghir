package github

import (
	"context"

	"github.com/google/go-github/v89/github"
)

func (c *Client) EditRelease(ctx context.Context, owner, repo string, id int64) error {
	// GraphQL API does not support updating a release.
	// https://docs.github.com/en/graphql/reference/mutations
	// https://pkg.go.dev/github.com/google/go-github/v89/github#RepositoriesService.UpdateRelease
	_, _, err := c.repos.UpdateRelease(ctx, owner, repo, id, github.UpdateReleaseRequest{})
	return err //nolint:wrapcheck
}
