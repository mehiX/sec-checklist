package application

import "github.com/mehix/sec-checklist/pkg/domain/application/db"

type Option func(*service)

func WithDb(dsn string) Option {
	return func(s *service) {
		s.dbRepo = db.NewRepository(dsn)
	}
}
