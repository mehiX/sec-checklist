package application

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mehix/sec-checklist/pkg/domain"
	"github.com/mehix/sec-checklist/pkg/domain/application"
	"github.com/mehix/sec-checklist/pkg/domain/check"
	"github.com/mehix/sec-checklist/pkg/iFacts"
)

var ErrDbNotConnected = fmt.Errorf("not connected to database")

type ApplicationService interface {
	FetchApplicationByID(context.Context, string) (*domain.Application, error)
	ListAllApplications(context.Context) ([]domain.Application, error)
	SaveApplication(context.Context, *domain.Application) error
	SaveApplicationFilters(context.Context, *domain.Application) error
	SaveApplicationOrImportFromIFacts(context.Context, *domain.Application, iFacts.Client) (*domain.Application, error)
	SaveFromIFacts(ctx context.Context, iFactsID string, ifclient iFacts.Client) error

	FilterControls(context.Context, domain.ControlsFilter) ([]domain.Control, error)
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

func (s *service) FilterControls(ctx context.Context, filter domain.ControlsFilter) ([]domain.Control, error) {
	return s.dbRepo.ControlsForFilter(ctx, &filter)
}

func (s *service) SaveApplicationFilters(ctx context.Context, app *domain.Application) error {
	if s.appsDbRepo == nil {
		return ErrDbNotConnected
	}

	return s.appsDbRepo.SaveFilters(ctx, app)
}

func (s *service) SaveApplicationOrImportFromIFacts(ctx context.Context, app *domain.Application, cli iFacts.Client) (*domain.Application, error) {
	if s.appsDbRepo == nil {
		return nil, ErrDbNotConnected
	}

	if app.ID == "" {
		found, err := s.appsDbRepo.FindByInternalID(ctx, app.InternalID)
		if err != nil {
			log.Printf("finding app by internal ID, got error: %v\n", err.Error())
		} else {
			if found != nil {
				// we already have an app saved for this InternalID
				return found, nil
			}
		}

		app.ID = uuid.NewString()

		if err := s.appsDbRepo.Save(ctx, app); err != nil {
			return app, err
		}
	}

	return app, nil
}

func (s *service) SaveFromIFacts(ctx context.Context, id string, ifc iFacts.Client) error {
	classifications, err := ifc.GetClassifications(id)
	if err != nil {
		return err
	}

	return s.appsDbRepo.SaveIFactsClassifications(ctx, id, classifications)
}
