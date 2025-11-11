package di

import (
	"github.com/r0manch1k/umbrella/signature-server/internal/repository"
	"github.com/r0manch1k/umbrella/signature-server/internal/repository/license"
)

type Repositories struct {
	License repository.LicenseRepository
}

func initRepositories(clients Clients) Repositories {
	licenseRepo := license.New(clients.PgSql)

	return Repositories{License: licenseRepo}
}
