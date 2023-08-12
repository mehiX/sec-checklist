package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mehix/sec-checklist/pkg/domain"
	"github.com/mehix/sec-checklist/pkg/domain/application"
	"github.com/mehix/sec-checklist/pkg/domain/check"
)

var ErrDbNotConnected = fmt.Errorf("not connected to database")

type ApplicationService interface {
	FetchApplicationByID(context.Context, string) (*domain.Application, error)
	ListAllApplications(context.Context) ([]domain.Application, error)
	SaveApplication(context.Context, *domain.Application) error
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

func (s *service) FetchApplicationByID(ctx context.Context, id string) (*domain.Application, error) {
	if s.appsDbRepo == nil {
		return nil, ErrDbNotConnected
	}

	return s.appsDbRepo.FetchByID(ctx, id)
}

func (s *service) ListAllApplications(ctx context.Context) ([]domain.Application, error) {
	if s.appsDbRepo == nil {
		return nil, ErrDbNotConnected
	}

	return s.appsDbRepo.ListAll(ctx)
}

func (s *service) SaveApplication(ctx context.Context, app *domain.Application) error {
	if s.appsDbRepo == nil {
		return ErrDbNotConnected
	}

	if app.ID == "" {
		app.ID = uuid.NewString()

		if err := s.appsDbRepo.Save(ctx, app); err != nil {
			return err
		}

		ctrls, err := s.dbRepo.ControlsForFilter(ctx, &domain.ControlsFilter{
			OnlyHandleCentrally:         &app.OnlyHandledCentrally,
			HandledCentrallyBy:          &app.HandledCentrallyBy,
			ExcludeForExternalSupplier:  &app.ExcludeForExternalSupplier,
			SoftwareDevelopmentRelevant: &app.SoftwareDevelopmentRelevant,
			CloudOnly:                   &app.CloudOnly,
			PhysicalSecurityOnly:        &app.PhysicalSecurityOnly,
			PersonalSecurityOnly:        &app.PersonalSecurityOnly,
		})
		if err != nil {
			return err
		}

		return s.dbRepo.SaveForApplication(ctx, app, ctrls)
	}

	return s.appsDbRepo.Update(ctx, app)
}
