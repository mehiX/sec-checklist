package db

import "github.com/mehix/sec-checklist/pkg/domain"

// Scanner used as constraint for `scanForControl`
type Scanner interface {
	Scan(...any) error
}

func scanForControl[T Scanner](s T) (domain.Control, error) {
	var id, tp, name, desc, c, i, a, t, handledCentrallyBy string
	var onlyHandleCentrally, excludedForExternalSupplier, softwareDevelopmentRelevant bool
	var cloudOnly, physicalSecurityOnly, personalSecurityOnly bool
	finalValues := []any{&id, &tp, &name, &desc, &c, &i, &a, &t,
		&onlyHandleCentrally, &handledCentrallyBy, &excludedForExternalSupplier,
		&softwareDevelopmentRelevant, &cloudOnly,
		&physicalSecurityOnly, &personalSecurityOnly}

	if err := s.Scan(finalValues...); err != nil {
		return domain.Control{}, err
	}

	return domain.Control{
		ID:                          id,
		Type:                        tp,
		Name:                        name,
		Description:                 desc,
		C:                           c,
		I:                           i,
		A:                           a,
		T:                           t,
		OnlyHandledCentrally:        onlyHandleCentrally,
		HandledCentrallyBy:          handledCentrallyBy,
		ExcludeForExternalSupplier:  excludedForExternalSupplier,
		SoftwareDevelopmentRelevant: softwareDevelopmentRelevant,
		CloudOnly:                   cloudOnly,
		PhysicalSecurityOnly:        physicalSecurityOnly,
		PersonalSecurityOnly:        personalSecurityOnly,
	}, nil
}
