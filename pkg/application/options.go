package application

import (
	"github.com/mehix/sec-checklist/pkg/domain/application/db"
	controlsDb "github.com/mehix/sec-checklist/pkg/domain/check/db"
	"github.com/mehix/sec-checklist/pkg/domain/check/xls"
)

type Option func(Service)

func WithAppsDb(dsn string) Option {
	return func(s Service) {
		s.(*service).appsDbRepo = db.NewRepository(dsn)
	}
}

func WithXls(fpath, sheetName string) Option {
	return func(s Service) {
		s.(*service).xlsRepo = xls.NewRepository(fpath, sheetName)
	}
}

func WithControlsDb(dsn string) Option {
	return func(s Service) {
		s.(*service).dbRepo = controlsDb.NewRepository(dsn)
	}
}
