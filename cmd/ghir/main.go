package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/suzuki-shunsuke/ghir/pkg/cli"
	"github.com/suzuki-shunsuke/ghir/pkg/log"
	"github.com/suzuki-shunsuke/go-stdutil"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logger, logLevel := log.New(os.Stderr, version)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := cli.Run(ctx, logger, logLevel, &stdutil.LDFlags{
		Version: version,
		Commit:  commit,
		Date:    date,
	}); err != nil {
		slogerr.WithError(logger, err).Error("ghir failed")
		return 1
	}
	return 0
}
