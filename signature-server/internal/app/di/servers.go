package di

import (
	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
	"github.com/r0manch1k/umbrella/signature-server/pkg/servers"
	"github.com/r0manch1k/umbrella/signature-server/pkg/servers/httpserver"
)

type Servers struct {
	HTTP servers.Server
}

func initServers(cfg *config.Config, l *logger.Logger, controllers Controllers) Servers {
	httpSrv := httpserver.New(
		l,
		controllers.HTTP.Router.Handler,
		cfg.App.Name,
		httpserver.Address(cfg.HTTP.Host, cfg.HTTP.Port),
	)

	return Servers{HTTP: httpSrv}
}
