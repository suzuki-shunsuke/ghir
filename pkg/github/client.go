package github

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/go-github/v81/github"
	"github.com/shurcooL/githubv4"
	"github.com/suzuki-shunsuke/ghtkn-go-sdk/ghtkn"
	"github.com/suzuki-shunsuke/go-retryablehttp"
	"golang.org/x/oauth2"
)

type Client struct {
	repos RepositoriesService
	v4    GraphQL
}

type InputNew struct {
	GHTKNEnabled bool
	AccessToken  string
}

func New(ctx context.Context, logger *slog.Logger, input *InputNew) (*Client, error) {
	httpClient, err := newHTTPClient(ctx, logger, input)
	if err != nil {
		return nil, err
	}
	return &Client{
		repos: github.NewClient(httpClient).Repositories,
		v4:    githubv4.NewClient(httpClient),
	}, nil
}

func newHTTPClient(ctx context.Context, logger *slog.Logger, input *InputNew) (*http.Client, error) {
	ts, err := newTokenSource(logger, input)
	if err != nil {
		return nil, err
	}
	return makeRetryable(oauth2.NewClient(ctx, ts), logger), nil
}

func newTokenSource(logger *slog.Logger, input *InputNew) (oauth2.TokenSource, error) {
	if input.GHTKNEnabled {
		client := ghtkn.New()
		return client.TokenSource(logger, &ghtkn.InputGet{}), nil
	}
	if input.AccessToken != "" {
		return oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: input.AccessToken},
		), nil
	}
	return nil, errors.New("either GHTKNEnabled or AccessToken must be set")
}

func makeRetryable(client *http.Client, logger *slog.Logger) *http.Client {
	c := retryablehttp.NewClient()
	c.HTTPClient = client
	c.Logger = logger
	return c.StandardClient()
}

type RepositoriesService interface {
	EditRelease(ctx context.Context, owner, repo string, id int64, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
}

type GraphQL interface {
	Query(ctx context.Context, q any, variables map[string]any) error
}
