package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/internal/app/di"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
	"github.com/r0manch1k/umbrella/signature-server/pkg/servers"
)

func Run(cfg *config.Config) {
	loc, err := time.LoadLocation(cfg.App.TZ)
	if err != nil {
		panic(fmt.Errorf("invalid timezone %s: %w", cfg.App.TZ, err))
	}

	time.Local = loc

	deps, err := di.New(cfg)
	if err != nil {
		panic(err)
	}
	defer deps.Clients.PgSQL.Close()

	deps.Logger.Info("%s", fmt.Sprintf("Starting %s app...", cfg.App.Name))

	deps.Servers.HTTP.Start()
	gracefulShutdown(deps.Logger, deps.Servers.HTTP)
}

func gracefulShutdown(l *logger.Logger, srvs ...servers.Server) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-interrupt:
		l.Info("%s", fmt.Sprintf("App - Run - signal: %s", sig))
	case err := <-mergeErrors(srvs...):
		l.Fatal(err)
	}

	for _, srv := range srvs {
		if err := srv.Shutdown(); err != nil {
			l.Fatal(err)
		}
	}
}

func mergeErrors(srvs ...servers.Server) <-chan error {
	out := make(chan error, len(srvs))
	for _, srv := range srvs {
		go func(s servers.Server) {
			for err := range s.Notify() {
				out <- err
			}
		}(srv)
	}

	return out
}
