package check

import (
	"fmt"
)

type Control struct {
	Type                        string
	ID                          string
	Name                        string
	Description                 string
	C                           string
	I                           string
	A                           string
	T                           string
	PD                          string
	NSI                         string
	SESE                        string
	OTCL                        string
	CSRDirection                string // CS&R direction for control type
	SPSA                        string
	SPSAUnique                  string
	GDPR                        bool
	GDPRUnique                  bool
	ExternalSupplier            bool //
	AssetType                   string
	OperationalCapability       string
	PartOfGISR                  bool
	LastUpdated                 string
	OldID                       string
	OnlyHandledCentrally        bool
	HandledCentrallyBy          string
	ExcludeForExternalSupplier  bool
	SoftwareDevelopmentRelevant bool
	CloudOnly                   bool
	PhysicalSecurityOnly        bool
	PersonalSecurityOnly        bool
}

func (e Control) String() string {
	return fmt.Sprintf("%s - %s [%s]", e.ID, e.Name, e.Type)
}

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
