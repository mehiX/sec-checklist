package check

import "context"

type Reader interface {
	FetchAll() ([]Control, error)
	FetchByType(string) ([]Control, error)
	FetchByID(context.Context, string) (Control, error)
}

type Writer interface {
	SaveAll(context.Context, []Control) error
}

type ReaderWriter interface {
	Reader
	Writer
}
