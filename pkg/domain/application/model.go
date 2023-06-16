package application

type Application struct {
	ID                          string
	Name                        string
	OnlyHandledCentrally        bool
	HandledCentrallyBy          string
	ExcludeForExternalSupplier  bool
	SoftwareDevelopmentRelevant bool
	CloudOnly                   bool
	PhysicalSecurityOnly        bool
	PersonalSecurityOnly        bool
}
