package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghir/pkg/github"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type InputRun struct {
	RepoOwner string
	RepoName  string
}

func (c *Controller) Run(ctx context.Context, logger *slog.Logger, input *InputRun) error {
	releases, err := c.input.GitHub.ListReleases(ctx, input.RepoOwner, input.RepoName)
	if err != nil {
		return fmt.Errorf("list releases: %w", err)
	}
	if len(releases) == 0 {
		logger.Info("no releases found")
		return nil
	}
	rs := make([]*github.Release, 0, len(releases))
	for _, release := range releases {
		logger := logger.With("tag", release.TagName, "release_id", release.DatabaseID)
		if release.Immutable {
			logger.Debug("ignore immutable release")
			continue
		}
		if release.IsDraft {
			logger.Debug("ignore draft release")
			continue
		}
		rs = append(rs, release)
	}
	if len(rs) == 0 {
		logger.Info("all releases are immutable")
		return nil
	}
	logger.Info("found mutable releases", "count", len(rs))
	for _, release := range rs {
		logger.Info("editing release description to make the release immutable", "tag", release.TagName, "release_id", release.DatabaseID)
		if err := c.input.GitHub.EditRelease(ctx, input.RepoOwner, input.RepoName, release.DatabaseID, release.Description+"\n"); err != nil {
			return fmt.Errorf("append a newline to the description of release: %w", slogerr.With(err, "tag", release.TagName, "release_id", release.DatabaseID))
		}
		if err := c.input.GitHub.EditRelease(ctx, input.RepoOwner, input.RepoName, release.DatabaseID, release.Description); err != nil {
			return fmt.Errorf("remove a newline from the description of release: %w", slogerr.With(err, "tag", release.TagName, "release_id", release.DatabaseID))
		}
	}
	return nil
}
