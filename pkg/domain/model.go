package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type Application struct {
	ID                          string           `json:"id"`
	InternalID                  int              `json:"app_internal_id"`
	Name                        string           `json:"app_name"`
	IFactsID                    string           `json:"ifacts_id"`
	OnlyHandledCentrally        bool             `json:"only_handle_centrally"`
	HandledCentrallyBy          string           `json:"handled_centrally_by"`
	ExcludeForExternalSupplier  bool             `json:"exclude_for_external_supplier"`
	SoftwareDevelopmentRelevant bool             `json:"software_development_relevant"`
	CloudOnly                   bool             `json:"cloud_only"`
	PhysicalSecurityOnly        bool             `json:"physical_security_only"`
	PersonalSecurityOnly        bool             `json:"personal_security_only"`
	Classifications             []Classification `json:"classifications"`
}

func (a *Application) LevelC() string {
	return getLevel("(C)", a.Classifications)
}
func (a *Application) LevelI() string {
	return getLevel("(I)", a.Classifications)
}
func (a *Application) LevelA() string {
	return getLevel("(A)", a.Classifications)
}
func (a *Application) LevelT() string {
	return getLevel("(T)", a.Classifications)
}

var num = regexp.MustCompile("[0-9]{1}")

// getLevel returns level as string since it is easier to use it
// like this when filtering controls.
func getLevel(name string, classifications []Classification) string {
	if classifications == nil {
		return ""
	}

	for _, c := range classifications {
		if strings.Contains(c.Name, name) {
			b := num.Find([]byte(c.LevelName))
			if b != nil {
				return string(b)
			}
		}
	}

	return ""

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

type Classification struct {
	ID        string `json:"ClassificationId"`
	Name      string `json:"ClassificationName"`
	LevelID   string `json:"SavedLevelId"`
	LevelName string `json:"SavedLevelName"`
}
