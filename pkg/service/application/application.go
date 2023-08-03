package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mehix/sec-checklist/pkg/domain/application"
	"github.com/mehix/sec-checklist/pkg/domain/check"
)

var ErrDbNotConnected = fmt.Errorf("not connected to database")

type ApplicationService interface {
	FetchApplicationByID(context.Context, string) (*application.Application, error)
	ListAll(context.Context) ([]application.Application, error)
	Save(context.Context, *application.Application) error
	Update(context.Context, *application.Application) error
}

type service struct {
	appsDbRepo application.ReaderWriter
	xlsRepo    check.Reader
	dbRepo     check.ReaderWriter
}

func NewService(options ...Option) Service {
	s := &service{}
	for _, o := range options {
		o(s)
	}
	return s
}

func (s *service) FetchApplicationByID(ctx context.Context, id string) (*application.Application, error) {
	if s.appsDbRepo == nil {
		return nil, ErrDbNotConnected
	}

	return s.appsDbRepo.FetchByID(ctx, id)
}

func (s *service) ListAll(ctx context.Context) ([]application.Application, error) {
	if s.appsDbRepo == nil {
		return nil, ErrDbNotConnected
	}

	return s.appsDbRepo.ListAll(ctx)
}

func (s *service) Save(ctx context.Context, app *application.Application) error {
	if s.appsDbRepo == nil {
		return ErrDbNotConnected
	}

	app.ID = uuid.NewString()

	if err := s.appsDbRepo.Save(ctx, app); err != nil {
		return err
	}

	ctrls, err := s.FetchByApplication(ctx, app)
	if err != nil {
		return err
	}

	return s.dbRepo.SaveForApplication(ctx, app, ctrls)
}

func (s *service) Update(ctx context.Context, app *application.Application) error {
	if s.appsDbRepo == nil {
		return ErrDbNotConnected
	}

	return s.appsDbRepo.Update(ctx, app)
}
