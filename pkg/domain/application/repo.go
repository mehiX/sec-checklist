package application

import "context"

type Reader interface {
	FetchAppByID(context.Context, string) (*Application, error)
	ListAllApps(context.Context) ([]Application, error)
}

type Writer interface {
	SaveApp(context.Context, *Application) error
}

type ReaderWriter interface {
	Reader
	Writer
}
