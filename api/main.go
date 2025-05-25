package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"os"
	"runtime"
	"service/foundation/logger"

	"github.com/ardanlabs/conf/v3"
)

var build = "develop"

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

	// -------------------------------------------------------------------------
	// Configuration
	cfg := struct {
		conf.Version
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "sales service",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for ouput %w", err)
	}

	log.Info(ctx, "startup", "config", out)

	log.BuildInfo(ctx)

	expvar.NewString("build").Set(cfg.Build)

	return nil
}
