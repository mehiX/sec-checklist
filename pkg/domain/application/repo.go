package application

import (
	"context"

	"github.com/mehix/sec-checklist/pkg/domain"
)

type Reader interface {
	FetchByID(context.Context, string) (*domain.Application, error)
	FindByInternalID(context.Context, int) (*domain.Application, error)
	ListAll(context.Context) ([]domain.Application, error)
}

type Writer interface {
	Save(context.Context, *domain.Application) error
	SaveFilters(context.Context, *domain.Application) error
	Update(context.Context, *domain.Application) error
}

type ReaderWriter interface {
	Reader
	Writer
}
