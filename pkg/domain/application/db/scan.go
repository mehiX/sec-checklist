package db

import (
	"github.com/mehix/sec-checklist/pkg/domain"
)

// Scanner used as constraint for `scanForControl`
type Scanner interface {
	Scan(...any) error
}

func scanForApp[T Scanner](s T) (*domain.Application, error) {

	var id, name, iFactsID, handledCentrallyBy string
	var internalID int
	var onlyHandleCentrally, excludedForExternalSupplier, softwareDevelopmentRelevant bool
	var cloudOnly, physicalSecurityOnly, personalSecurityOnly bool

	finalValues := []any{
		&id, &internalID, &name, &iFactsID,
		&onlyHandleCentrally, &handledCentrallyBy, &excludedForExternalSupplier,
		&softwareDevelopmentRelevant, &cloudOnly,
		&physicalSecurityOnly, &personalSecurityOnly,
	}

	if err := s.Scan(finalValues...); err != nil {
		return nil, err
	}

	return &domain.Application{
		ID:                          id,
		Name:                        name,
		IFactsID:                    iFactsID,
		InternalID:                  internalID,
		OnlyHandledCentrally:        onlyHandleCentrally,
		HandledCentrallyBy:          handledCentrallyBy,
		ExcludeForExternalSupplier:  excludedForExternalSupplier,
		SoftwareDevelopmentRelevant: softwareDevelopmentRelevant,
		CloudOnly:                   cloudOnly,
		PhysicalSecurityOnly:        physicalSecurityOnly,
		PersonalSecurityOnly:        personalSecurityOnly,
	}, nil
}
