package check

import (
	"context"

	"github.com/mehix/sec-checklist/pkg/domain/application"
)

type Reader interface {
	FetchAll() ([]Control, error)
	FetchByType(string) ([]Control, error)
	FetchByID(context.Context, string) (Control, error)
	FetchForApplication(context.Context, *application.Application) ([]Control, error)
}

type Writer interface {
	SaveAll(context.Context, []Control) error
}

type ReaderWriter interface {
	Reader
	Writer
}
