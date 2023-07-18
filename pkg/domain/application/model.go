package application

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
