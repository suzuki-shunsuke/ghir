//nolint:funlen
package controller_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/suzuki-shunsuke/ghir/pkg/controller"
	"github.com/suzuki-shunsuke/ghir/pkg/github"
)

type mockGitHub struct {
	listReleasesFunc func(ctx context.Context, owner, repo string) ([]*github.Release, error)
	editReleaseFunc  func(ctx context.Context, owner, repo string, releaseID int64) error
}

func (m *mockGitHub) ListReleases(ctx context.Context, owner, repo string) ([]*github.Release, error) {
	if m.listReleasesFunc != nil {
		return m.listReleasesFunc(ctx, owner, repo)
	}
	return nil, nil
}

func (m *mockGitHub) EditRelease(ctx context.Context, owner, repo string, releaseID int64) error {
	if m.editReleaseFunc != nil {
		return m.editReleaseFunc(ctx, owner, repo, releaseID)
	}
	return nil
}

func TestController_Run(t *testing.T) {
	t.Parallel()
	logger := slog.Default()

	tests := []struct {
		name           string
		input          *controller.InputRun
		mockGitHub     *mockGitHub
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "no releases found",
			input: &controller.InputRun{
				RepoOwner: "owner",
				RepoName:  "repo",
			},
			mockGitHub: &mockGitHub{
				listReleasesFunc: func(_ context.Context, _, _ string) ([]*github.Release, error) {
					return []*github.Release{}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "all releases are immutable",
			input: &controller.InputRun{
				RepoOwner: "owner",
				RepoName:  "repo",
			},
			mockGitHub: &mockGitHub{
				listReleasesFunc: func(_ context.Context, _, _ string) ([]*github.Release, error) {
					return []*github.Release{
						{
							TagName:    "v1.0.0",
							DatabaseID: 123,
							Immutable:  true,
							IsDraft:    false,
						},
					}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "all releases are drafts",
			input: &controller.InputRun{
				RepoOwner: "owner",
				RepoName:  "repo",
			},
			mockGitHub: &mockGitHub{
				listReleasesFunc: func(_ context.Context, _, _ string) ([]*github.Release, error) {
					return []*github.Release{
						{
							TagName:    "v1.0.0",
							DatabaseID: 123,
							Immutable:  false,
							IsDraft:    true,
						},
					}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "successfully edit mutable release",
			input: &controller.InputRun{
				RepoOwner: "owner",
				RepoName:  "repo",
			},
			mockGitHub: &mockGitHub{
				listReleasesFunc: func(_ context.Context, _, _ string) ([]*github.Release, error) {
					return []*github.Release{
						{
							TagName:    "v1.0.0",
							DatabaseID: 123,
							Immutable:  false,
							IsDraft:    false,
						},
					}, nil
				},
				editReleaseFunc: func(_ context.Context, _, _ string, releaseID int64) error {
					// This is a simplified check - in a real test you might want to track call order
					if releaseID != 123 {
						t.Errorf("unexpected releaseID: got %d, want 123", releaseID)
					}
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "list releases error",
			input: &controller.InputRun{
				RepoOwner: "owner",
				RepoName:  "repo",
			},
			mockGitHub: &mockGitHub{
				listReleasesFunc: func(_ context.Context, _, _ string) ([]*github.Release, error) {
					return nil, errors.New("list releases failed")
				},
			},
			wantErr:        true,
			wantErrMessage: "list releases: list releases failed",
		},
		{
			name: "edit release error",
			input: &controller.InputRun{
				RepoOwner: "owner",
				RepoName:  "repo",
			},
			mockGitHub: &mockGitHub{
				listReleasesFunc: func(_ context.Context, _, _ string) ([]*github.Release, error) {
					return []*github.Release{
						{
							TagName:    "v1.0.0",
							DatabaseID: 123,
							Immutable:  false,
							IsDraft:    false,
						},
					}, nil
				},
				editReleaseFunc: func(_ context.Context, _, _ string, _ int64) error {
					return errors.New("edit release failed")
				},
			},
			wantErr:        true,
			wantErrMessage: "update the release: edit release failed",
		},
		{
			name: "multiple mutable releases",
			input: &controller.InputRun{
				RepoOwner: "owner",
				RepoName:  "repo",
			},
			mockGitHub: &mockGitHub{
				listReleasesFunc: func(_ context.Context, _, _ string) ([]*github.Release, error) {
					return []*github.Release{
						{
							TagName:    "v1.0.0",
							DatabaseID: 123,
							Immutable:  false,
							IsDraft:    false,
						},
						{
							TagName:    "v2.0.0",
							DatabaseID: 456,
							Immutable:  true,
							IsDraft:    false,
						},
						{
							TagName:    "v3.0.0",
							DatabaseID: 789,
							Immutable:  false,
							IsDraft:    false,
						},
					}, nil
				},
				editReleaseFunc: func(_ context.Context, _, _ string, releaseID int64) error {
					// Should only be called for releases 123 and 789, not 456
					if releaseID != 123 && releaseID != 789 {
						t.Errorf("unexpected releaseID: got %d, should only edit 123 or 789", releaseID)
					}
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := controller.New(&controller.Input{
				GitHub: tt.mockGitHub,
			})

			if err := ctrl.Run(t.Context(), logger, tt.input); err != nil {
				if !tt.wantErr {
					t.Error(err)
					return
				}
				if tt.wantErrMessage != "" && err.Error() != tt.wantErrMessage {
					t.Errorf("Controller.Run() error = %v, wantErrMessage %v", err.Error(), tt.wantErrMessage)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("Controller.Run() error = nil, wantErr %v", tt.wantErr)
				return
			}
		})
	}
}

func TestController_Run_EditReleaseCallOrder(t *testing.T) {
	t.Parallel()
	logger := slog.Default()

	type call struct {
		ReleaseID int64
	}

	var calls []call

	mockGitHub := &mockGitHub{
		listReleasesFunc: func(_ context.Context, _, _ string) ([]*github.Release, error) {
			return []*github.Release{
				{
					TagName:    "v1.0.0",
					DatabaseID: 123,
					Immutable:  false,
					IsDraft:    false,
				},
			}, nil
		},
		editReleaseFunc: func(_ context.Context, _, _ string, releaseID int64) error {
			calls = append(calls, call{
				ReleaseID: releaseID,
			})
			return nil
		},
	}

	ctrl := controller.New(&controller.Input{
		GitHub: mockGitHub,
	})

	err := ctrl.Run(t.Context(), logger, &controller.InputRun{
		RepoOwner: "owner",
		RepoName:  "repo",
	})
	if err != nil {
		t.Fatalf("Controller.Run() error = %v", err)
	}

	expectedCalls := []call{
		{
			ReleaseID: 123,
		},
	}

	if diff := cmp.Diff(expectedCalls, calls); diff != "" {
		t.Errorf("EditRelease calls mismatch (-want +got):\n%s", diff)
	}
}
