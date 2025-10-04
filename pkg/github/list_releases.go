package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

/*
query($owner: String!, $repo: String!) {
  repository(owner: $owner, name: $repo) {
    releases(first:100, after:"") {
      pageInfo {
        hasNextPage
        endCursor
      }
      nodes {
        id
        isDraft
        immutable
        tagName
        description
      }
    }
  }
}
*/

type ListReleasesQuery struct {
	Repository *Repository `graphql:"repository(owner: $repoOwner, name: $repoName)"`
}

func (q *ListReleasesQuery) PageInfo() *PageInfo {
	return q.Repository.Releases.PageInfo
}

func (q *ListReleasesQuery) Nodes() []*Release {
	return q.Repository.Releases.Nodes
}

type Repository struct {
	Releases *Releases `graphql:"releases(first:100, after:$cursor)"`
}

type Releases struct {
	PageInfo *PageInfo  `json:"pageInfo"`
	Nodes    []*Release `json:"nodes"`
}

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type Release struct {
	Description string `json:"description"`
	Immutable   bool   `json:"immutable"`
	IsDraft     bool   `json:"isDraft"`
	TagName     string `json:"tagName"`
	DatabaseID  int64  `json:"databaseId"`
}

func (c *Client) ListReleases(ctx context.Context, owner, repo string) ([]*Release, error) {
	// https://docs.github.com/en/graphql/reference/objects#release
	var releases []*Release
	var cursor string
	variables := map[string]any{
		"repoOwner": githubv4.String(owner),
		"repoName":  githubv4.String(repo),
		"cursor":    githubv4.String(cursor),
	}
	for range 100 {
		q := &ListReleasesQuery{}
		if err := c.v4.Query(ctx, q, variables); err != nil {
			return nil, fmt.Errorf("list releases by GitHub GraphQL API: %w", err)
		}
		releases = append(releases, q.Nodes()...)
		pageInfo := q.PageInfo()
		if !pageInfo.HasNextPage {
			return releases, nil
		}
		variables["cursor"] = githubv4.String(pageInfo.EndCursor)
	}
	return releases, nil
}
