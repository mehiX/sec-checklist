package application

import (
	"context"
	"fmt"

	"github.com/mehix/sec-checklist/pkg/domain"
)

var ErrNoDb = fmt.Errorf("not connected to a database")
var ErrNoExcel = fmt.Errorf("no Excel file specified")

type ControlsService interface {
	FetchControlsFromExcel() ([]domain.Control, error)
	FetchAllControls() ([]domain.Control, error)
	FetchControlsByType(string) ([]domain.Control, error)
	FetchControlByID(context.Context, string) (domain.Control, error)
	SaveAllControls(context.Context, []domain.Control) error

	FetchAppControlByID(context.Context, *domain.Application, string) (domain.AppControl, error)
	FetchControlsByApplication(context.Context, *domain.Application) ([]domain.AppControl, error)
	SaveAppControl(context.Context, *domain.AppControl) error
}

func (s *service) FetchControlsFromExcel() ([]domain.Control, error) {
	if s.xlsRepo == nil {
		return nil, ErrNoExcel
	}
	return s.xlsRepo.FetchAll()
}

func (s *service) FetchAllControls() ([]domain.Control, error) {
	if s.dbRepo == nil {
		return nil, ErrNoDb
	}
	return s.dbRepo.FetchAll()
}

func (s *service) FetchControlsByType(t string) ([]domain.Control, error) {
	if s.dbRepo == nil {
		return nil, ErrNoDb
	}
	return s.dbRepo.FetchByType(t)
}

func (s *service) SaveAllControls(ctx context.Context, all []domain.Control) error {
	if s.dbRepo == nil {
		return ErrNoDb
	}
	return s.dbRepo.SaveAll(ctx, all)
}

func (s *service) FetchControlByID(ctx context.Context, id string) (domain.Control, error) {
	if s.dbRepo == nil {
		return domain.Control{}, ErrNoDb
	}
	return s.dbRepo.FetchByID(ctx, id)
}

func (s *service) FetchControlsByApplication(ctx context.Context, app *domain.Application) ([]domain.AppControl, error) {
	if s.dbRepo == nil {
		return nil, ErrNoDb
	}

	return s.dbRepo.ControlsForApplication(ctx, app.ID)
}

func (s *service) FetchAppControlByID(ctx context.Context, app *domain.Application, id string) (domain.AppControl, error) {
	if s.dbRepo == nil {
		return domain.AppControl{}, ErrNoDb
	}

	return s.dbRepo.ControlForApplicationByID(ctx, app.ID, id)
}

func (s *service) SaveAppControl(ctx context.Context, ctrl *domain.AppControl) error {
	if s.dbRepo == nil {
		return ErrNoDb
	}

	return s.dbRepo.SaveAppControl(ctx, ctrl)
}
