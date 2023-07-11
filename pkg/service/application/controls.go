package application

import (
	"context"
	"fmt"

	"github.com/mehix/sec-checklist/pkg/domain/application"
	"github.com/mehix/sec-checklist/pkg/domain/check"
)

var ErrNoDb = fmt.Errorf("not connected to a database")
var ErrNoExcel = fmt.Errorf("no Excel file specified")

type ControlsService interface {
	FetchAllFromExcel() ([]check.Control, error)
	FetchAll() ([]check.Control, error)
	FetchByType(string) ([]check.Control, error)
	FetchControlByID(context.Context, string) (check.Control, error)
	SaveAll(context.Context, []check.Control) error

	FetchByApplication(context.Context, *application.Application) ([]check.Control, error)
}

func (s *service) FetchAllFromExcel() ([]check.Control, error) {
	if s.xlsRepo == nil {
		return nil, ErrNoExcel
	}
	return s.xlsRepo.FetchAll()
}

func (s *service) FetchAll() ([]check.Control, error) {
	if s.dbRepo == nil {
		return nil, ErrNoDb
	}
	return s.dbRepo.FetchAll()
}

func (s *service) FetchByType(t string) ([]check.Control, error) {
	if s.dbRepo == nil {
		return nil, ErrNoDb
	}
	return s.dbRepo.FetchByType(t)
}

func (s *service) SaveAll(ctx context.Context, all []check.Control) error {
	if s.dbRepo == nil {
		return ErrNoDb
	}
	return s.dbRepo.SaveAll(ctx, all)
}

func (s *service) FetchControlByID(ctx context.Context, id string) (check.Control, error) {
	if s.dbRepo == nil {
		return check.Control{}, ErrNoDb
	}
	return s.dbRepo.FetchByID(ctx, id)
}

func (s *service) FetchByApplication(ctx context.Context, app *application.Application) ([]check.Control, error) {
	if s.dbRepo == nil {
		return nil, ErrNoDb
	}

	return s.dbRepo.FetchForApplication(ctx, app)
}
