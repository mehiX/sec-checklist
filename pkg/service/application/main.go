package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mehix/sec-checklist/pkg/domain/application"
)

var ErrDbNotConnected = fmt.Errorf("not connected to database")

type Service interface {
	FetchByID(context.Context, string) (*application.Application, error)
	ListAll(context.Context) ([]application.Application, error)
	Save(context.Context, *application.Application) error
	Update(context.Context, *application.Application) error
}

type service struct {
	dbRepo application.ReaderWriter
}

func NewService(options ...Option) Service {
	s := &service{}
	for _, o := range options {
		o(s)
	}
	return s
}

func (s *service) FetchByID(ctx context.Context, id string) (*application.Application, error) {
	if s.dbRepo == nil {
		return nil, ErrDbNotConnected
	}

	return s.dbRepo.FetchByID(ctx, id)
}

func (s *service) ListAll(ctx context.Context) ([]application.Application, error) {
	if s.dbRepo == nil {
		return nil, ErrDbNotConnected
	}

	return s.dbRepo.ListAll(ctx)
}

func (s *service) Save(ctx context.Context, app *application.Application) error {
	if s.dbRepo == nil {
		return ErrDbNotConnected
	}

	app.ID = uuid.NewString()

	return s.dbRepo.Save(ctx, app)
}

func (s *service) Update(ctx context.Context, app *application.Application) error {
	if s.dbRepo == nil {
		return ErrDbNotConnected
	}

	return s.dbRepo.Update(ctx, app)
}
