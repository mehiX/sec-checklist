package db

import "github.com/mehix/sec-checklist/pkg/domain/application"

// Scanner used as constraint for `scanForControl`
type Scanner interface {
	Scan(...any) error
}

func scanForApp[T Scanner](s T) (*application.Application, error) {

	var id, name string
	var onlyHandleCentrally, excludedForExternalSupplier, softwareDevelopmentRelevant bool
	var cloudOnly, physicalSecurityOnly, personalSecurityOnly bool

	finalValues := []any{
		&id, &name,
		&onlyHandleCentrally, &excludedForExternalSupplier, &softwareDevelopmentRelevant,
		&cloudOnly, &physicalSecurityOnly, &personalSecurityOnly,
	}

	if err := s.Scan(finalValues...); err != nil {
		return nil, err
	}

	return &application.Application{
		ID:                          id,
		Name:                        name,
		OnlyHandledCentrally:        onlyHandleCentrally,
		ExcludeForExternalSupplier:  excludedForExternalSupplier,
		SoftwareDevelopmentRelevant: softwareDevelopmentRelevant,
		CloudOnly:                   cloudOnly,
		PhysicalSecurityOnly:        physicalSecurityOnly,
		PersonalSecurityOnly:        personalSecurityOnly,
	}, nil
}
