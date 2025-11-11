package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/internal/app/di"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
	"github.com/r0manch1k/umbrella/signature-server/pkg/servers"
)

func Run(cfg *config.Config) {
	deps, err := di.New(cfg)
	if err != nil {
		panic(err)
	}
	defer deps.Clients.PgSql.Close()

	deps.Logger.Info("%s", fmt.Sprintf("starting %s app...", cfg.App.Name))

	// Генерация ключей при необходимости
	if !fileExists(cfg.Signature.PrivateKeyPath) || !fileExists(cfg.Signature.PublicKeyPath) {
		deps.Logger.Info("keypair not found — creating new keys...")

		if err := deps.UseCases.KeyPair.GenerateAndSaveKeyPair(); err != nil {
			panic(err)
		}
	}

	deps.Servers.HTTP.Start()
	gracefulShutdown(deps.Logger, deps.Servers.HTTP)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func gracefulShutdown(l *logger.Logger, srvs ...servers.Server) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-interrupt:
		l.Info("%s", fmt.Sprintf("app - Run - signal: %s", sig))
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
