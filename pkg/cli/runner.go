package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"github.com/suzuki-shunsuke/ghir/pkg/controller"
	"github.com/suzuki-shunsuke/ghir/pkg/github"
	"github.com/suzuki-shunsuke/ghir/pkg/log"
	"github.com/suzuki-shunsuke/go-stdutil"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

var errInvalidRepoArg = errors.New("invalid repository name format")

const help = `ghir - Make GitHub Releases immutable
https://github.com/suzuki-shunsuke/ghir

USAGE:
   ghir --help [-h] # Show this help
   ghir --version [-v] # Show version
   ghir [--log-level <debug|info|warn|error>] [--enable-ghtkn] <owner>/<repo>

VERSION:
   %s
`

func Run(ctx context.Context, logger *slog.Logger, logLevel *slog.LevelVar, ldFlags *stdutil.LDFlags) error {
	// GITHUB_TOKEN, GHIR_GITHUB_TOKEN
	// GHIR_LOG_LEVEL -log-level
	// GHIR_ENABLE_GHTKN -enable-ghtkn
	flag := &Flag{}
	parseFlags(flag)

	if flag.Help {
		fmt.Fprintf(os.Stderr, help, ldFlags.Version)
		return nil
	}
	if flag.Version {
		fmt.Fprintf(os.Stderr, "%s\n", ldFlags.Version)
		return nil
	}

	if err := setLogLevel(logLevel, flag); err != nil {
		return err
	}

	if err := setEnableGHTKN(flag); err != nil {
		return err
	}

	repoOwner, repoName, err := validateRepo(flag.Args)
	if err != nil {
		return err
	}

	gh, err := github.New(ctx, logger, &github.InputNew{
		GHTKNEnabled: flag.EnableGHTKN,
		AccessToken:  getGitHubToken(),
	})
	if err != nil {
		return fmt.Errorf("create GitHub client: %w", err)
	}

	ctrl := controller.New(&controller.Input{
		GitHub: gh,
	})
	if err := ctrl.Run(ctx, logger, &controller.InputRun{
		RepoOwner: repoOwner,
		RepoName:  repoName,
	}); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

func getGitHubToken() string {
	if token := os.Getenv("GHIR_GITHUB_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITHUB_TOKEN")
}

func validateRepo(args []string) (string, string, error) {
	if len(args) == 0 {
		return "", "", errors.New("repo argument is required")
	}
	repo := args[0]
	repoOwner, repoName, ok := strings.Cut(repo, "/")
	if !ok {
		return "", "", slogerr.With(errInvalidRepoArg, "repository", repo) //nolint:wrapcheck
	}
	if strings.Contains(repoName, "/") {
		return "", "", slogerr.With(errInvalidRepoArg, "repository", repo) //nolint:wrapcheck
	}
	return repoOwner, repoName, nil
}

func setLogLevel(logLevel *slog.LevelVar, flag *Flag) error {
	if flag.LogLevel == "" {
		flag.LogLevel = os.Getenv(envLogLevel)
	}
	if flag.LogLevel != "" {
		if err := log.SetLevel(logLevel, flag.LogLevel); err != nil {
			return err //nolint:wrapcheck
		}
	}
	return nil
}

func setEnableGHTKN(flag *Flag) error {
	if flag.EnableGHTKN {
		return nil
	}
	s := os.Getenv("GHIR_ENABLE_GHTKN")
	b, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("GHIR_ENABLE_GHTKN must be boolean: %w", err)
	}
	flag.EnableGHTKN = b
	return nil
}

type Flag struct {
	LogLevel    string
	EnableGHTKN bool
	Help        bool
	Version     bool
	Args        []string
}

const envLogLevel = "GHIR_LOG_LEVEL"

func parseFlags(f *Flag) {
	pflag.StringVar(&f.LogLevel, "log-level", "", "log level (debug, info, warn, error)")
	pflag.BoolVar(&f.EnableGHTKN, "enable-ghtkn", false, "enable the integration with ghtkn")
	pflag.BoolVarP(&f.Help, "help", "h", false, "show help")
	pflag.BoolVarP(&f.Version, "version", "v", false, "show version")
	pflag.Parse()
	f.Args = pflag.Args()
}
