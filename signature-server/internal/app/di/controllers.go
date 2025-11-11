package di

import (
	"github.com/r0manch1k/umbrella/signature-server/internal/controller/http"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
)

type Controllers struct {
	HTTP *http.Router
}

func initControllers(l *logger.Logger, useCases UseCases) Controllers {
	httpRt := http.NewRouter(l, useCases.License, useCases.KeyPair)

	return Controllers{HTTP: httpRt}
}
