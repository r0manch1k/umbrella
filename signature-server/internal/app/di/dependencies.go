package di

import (
	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
)

type Dependencies struct {
	Clients      Clients
	Repositories Repositories
	Services     Services
	UseCases     UseCases
	Controllers  Controllers
	Servers      Servers
	Logger       *logger.Logger
}

func New(cfg *config.Config) (Dependencies, error) {
	l := logger.New(cfg.Log.Level, cfg.App.TZ)

	clients, err := initClients(cfg)
	if err != nil {
		return Dependencies{}, err
	}

	repositories := initRepositories(clients)

	services, err := initServices(cfg, repositories)
	if err != nil {
		return Dependencies{}, err
	}

	useCases := initUseCases(services)
	controllers := initControllers(l, useCases)
	servers := initServers(cfg, l, controllers)

	return Dependencies{
		Clients:      clients,
		Repositories: repositories,
		Services:     services,
		UseCases:     useCases,
		Controllers:  controllers,
		Servers:      servers,
		Logger:       l,
	}, nil
}
