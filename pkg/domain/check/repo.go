package check

import (
	"context"

	"github.com/mehix/sec-checklist/pkg/domain"
)

type Reader interface {
	FetchAll() ([]domain.Control, error)
	FetchByType(string) ([]domain.Control, error)
	FetchByID(context.Context, string) (domain.Control, error)
	ControlsForFilter(context.Context, *domain.ControlsFilter) ([]domain.Control, error)
	ControlsForApplication(ctx context.Context, appID string) ([]domain.AppControl, error)
	ControlForApplicationByID(ctx context.Context, appID string, ctrlID string) (domain.AppControl, error)
}

type Writer interface {
	SaveAll(context.Context, []domain.Control) error
	SaveForApplication(context.Context, *domain.Application, []domain.Control) error
	SaveAppControl(context.Context, *domain.AppControl) error
}

type ReaderWriter interface {
	Reader
	Writer
}
