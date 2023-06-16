package checks

import (
	"context"
	"fmt"

	"github.com/mehix/sec-checklist/pkg/domain/check"
)

var ErrNoDb = fmt.Errorf("not connected to a database")
var ErrNoExcel = fmt.Errorf("no Excel file specified")

type Service interface {
	FetchAllFromExcel() ([]check.Control, error)
	FetchAll() ([]check.Control, error)
	FetchByType(string) ([]check.Control, error)
	FetchByID(context.Context, string) (check.Control, error)
	SaveAll(context.Context, []check.Control) error
}

type service struct {
	xlsRepo check.Reader
	dbRepo  check.ReaderWriter
}

func NewService(options ...Option) Service {
	s := &service{}
	for _, o := range options {
		o(s)
	}
	return s
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

func (s *service) FetchByID(ctx context.Context, id string) (check.Control, error) {
	if s.dbRepo == nil {
		return check.Control{}, ErrNoDb
	}
	return s.dbRepo.FetchByID(ctx, id)
}
