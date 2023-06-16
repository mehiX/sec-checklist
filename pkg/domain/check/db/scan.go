package db

import "github.com/mehix/sec-checklist/pkg/domain/check"

// Scanner used as constraint for `scanForControl`
type Scanner interface {
	Scan(...any) error
}

func scanForControl[T Scanner](s T) (check.Control, error) {
	var id, tp, name, desc, assetType, lastUpdate, oldID, c, i, a, t, handledCentrallyBy string
	var pd, nsi, sese, otcl, csr, spsa, spsaUnique, operationalCapability string
	var gdpr, gdprUnique, externalSupplier, partOfGisr bool
	var onlyHandleCentrally, excludedForExternalSupplier, softwareDevelopmentRelevant bool
	var cloudOnly, physicalSecurityOnly, personalSecurityOnly bool
	finalValues := []any{
		&id, &tp, &name, &desc,
		&assetType, &lastUpdate, &oldID,
		&c, &i, &a, &t,
		&pd, &nsi, &sese, &otcl, &csr, &spsa, &spsaUnique,
		&gdpr, &gdprUnique, &externalSupplier, &operationalCapability, &partOfGisr,
		&onlyHandleCentrally, &handledCentrallyBy, &excludedForExternalSupplier,
		&softwareDevelopmentRelevant, &cloudOnly,
		&physicalSecurityOnly, &personalSecurityOnly}

	if err := s.Scan(finalValues...); err != nil {
		return check.Control{}, err
	}

	return check.Control{
		ID:                          id,
		Type:                        tp,
		Name:                        name,
		Description:                 desc,
		AssetType:                   assetType,
		LastUpdated:                 lastUpdate,
		OldID:                       oldID,
		C:                           c,
		I:                           i,
		A:                           a,
		T:                           t,
		PD:                          pd,
		NSI:                         nsi,
		SESE:                        sese,
		OTCL:                        otcl,
		CSRDirection:                csr,
		SPSA:                        spsa,
		SPSAUnique:                  spsaUnique,
		GDPR:                        gdpr,
		GDPRUnique:                  gdprUnique,
		ExternalSupplier:            externalSupplier,
		OperationalCapability:       operationalCapability,
		PartOfGISR:                  partOfGisr,
		OnlyHandledCentrally:        onlyHandleCentrally,
		HandledCentrallyBy:          handledCentrallyBy,
		ExcludeForExternalSupplier:  excludedForExternalSupplier,
		SoftwareDevelopmentRelevant: softwareDevelopmentRelevant,
		CloudOnly:                   cloudOnly,
		PhysicalSecurityOnly:        physicalSecurityOnly,
		PersonalSecurityOnly:        personalSecurityOnly,
	}, nil
}
