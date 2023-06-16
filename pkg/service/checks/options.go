package checks

import (
	"github.com/mehix/sec-checklist/pkg/domain/check/db"
	"github.com/mehix/sec-checklist/pkg/domain/check/xls"
)

type Option func(Service)

func WithXls(fpath, sheetName string) Option {
	return func(s Service) {
		s.(*service).xlsRepo = xls.NewRepository(fpath, sheetName)
	}
}

func WithDb(dsn string) Option {
	return func(s Service) {
		s.(*service).dbRepo = db.NewRepository(dsn)
	}
}
