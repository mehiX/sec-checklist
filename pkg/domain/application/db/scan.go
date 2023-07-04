package db

import "github.com/mehix/sec-checklist/pkg/domain/application"

// Scanner used as constraint for `scanForControl`
type Scanner interface {
	Scan(...any) error
}

func scanForApp[T Scanner](s T) (*application.Application, error) {

	var id, name, handledCentrallyBy string
	var internalID int
	var onlyHandleCentrally, excludedForExternalSupplier, softwareDevelopmentRelevant bool
	var cloudOnly, physicalSecurityOnly, personalSecurityOnly bool

	finalValues := []any{
		&id, &internalID, &name,
		&onlyHandleCentrally, &handledCentrallyBy, &excludedForExternalSupplier,
		&softwareDevelopmentRelevant, &cloudOnly,
		&physicalSecurityOnly, &personalSecurityOnly,
	}

	if err := s.Scan(finalValues...); err != nil {
		return nil, err
	}

	return &application.Application{
		ID:                          id,
		Name:                        name,
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
