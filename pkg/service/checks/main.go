package checks

import (
	"context"
	"fmt"

	"github.com/mehix/sec-checklist/pkg/domain"
)

type Service interface {
	FetchAllFromExcel() ([]domain.Control, error)
	FetchAll() ([]domain.Control, error)
	FetchByType(string) ([]domain.Control, error)
	FetchByID(context.Context, string) (domain.Control, error)
	SaveAll(context.Context, []domain.Control) error
}

type service struct {
	xlsRepo domain.Reader
	dbRepo  domain.ReaderWriter
}

func NewService(options ...Option) Service {
	s := &service{}
	for _, o := range options {
		o(s)
	}
	return s
}

func (s *service) FetchAllFromExcel() ([]domain.Control, error) {
	if s.xlsRepo == nil {
		return nil, fmt.Errorf("no Excel file was specified")
	}
	return s.xlsRepo.FetchAll()
}

func (s *service) FetchAll() ([]domain.Control, error) {
	return s.dbRepo.FetchAll()
}

func (s *service) FetchByType(t string) ([]domain.Control, error) {
	return s.dbRepo.FetchByType(t)
}

func (s *service) SaveAll(ctx context.Context, all []domain.Control) error {
	return s.dbRepo.SaveAll(ctx, all)
}

func (s *service) FetchByID(ctx context.Context, id string) (domain.Control, error) {
	return s.dbRepo.FetchByID(ctx, id)
}
