package domain

import "fmt"

type Application struct {
	ID                          string `json:"id"`
	InternalID                  int    `json:"app_internal_id"`
	Name                        string `json:"app_name"`
	OnlyHandledCentrally        bool   `json:"only_handle_centrally"`
	HandledCentrallyBy          string `json:"handled_centrally_by"`
	ExcludeForExternalSupplier  bool   `json:"exclude_for_external_supplier"`
	SoftwareDevelopmentRelevant bool   `json:"software_development_relevant"`
	CloudOnly                   bool   `json:"cloud_only"`
	PhysicalSecurityOnly        bool   `json:"physical_security_only"`
	PersonalSecurityOnly        bool   `json:"personal_security_only"`
	C                           int    `json:"classif_c"`
	I                           int    `json:"classif_i"`
	A                           int    `json:"classif_a"`
	T                           int    `json:"classif_t"`
}

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

type ControlsFilter struct {
	OnlyHandleCentrally         *bool   `json:"only_handle_centrally,omitempty"`
	HandledCentrallyBy          *string `json:"handled_centrally_by,omitempty"`
	ExcludeForExternalSupplier  *bool   `json:"exclude_for_external_supplier,omitempty"`
	SoftwareDevelopmentRelevant *bool   `json:"software_development_relevant,omitempty"`
	CloudOnly                   *bool   `json:"cloud_only,omitempty"`
	PhysicalSecurityOnly        *bool   `json:"physical_security_only,omitempty"`
	PersonalSecurityOnly        *bool   `json:"personal_security_only,omitempty"`
}

type AppControl struct {
	AppID       string
	ControlID   string
	Name        string
	Description string
	IsDone      bool
	Notes       string
}
