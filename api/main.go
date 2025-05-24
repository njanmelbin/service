package main

import (
	"context"
	"os"
	"runtime"
	"service/foundation/logger"
)

func main() {
	var log *logger.Logger

	ctx := context.Background()

	traceIDFn := func(ctx context.Context) string {
		return ""
	}
	log = logger.New(os.Stdout, logger.LevelInfo, "SALES", traceIDFn)

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	return nil
}
