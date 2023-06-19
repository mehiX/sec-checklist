package application

import "context"

type Reader interface {
	FetchByID(context.Context, string) (*Application, error)
	ListAll(context.Context) ([]Application, error)
}

type Writer interface {
	Save(context.Context, *Application) error
	Update(context.Context, *Application) error
}

type ReaderWriter interface {
	Reader
	Writer
}
