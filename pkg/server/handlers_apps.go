package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mehix/sec-checklist/pkg/domain"
	appDomain "github.com/mehix/sec-checklist/pkg/domain"
	"github.com/mehix/sec-checklist/pkg/iFacts"
	"github.com/mehix/sec-checklist/pkg/service/application"
)

func listAllApps(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		apps, err := svc.ListAllApplications(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(apps); err != nil {
			handleError(w, err)
			return
		}
	}
}

func showAppByID(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		app, ok := r.Context().Value(ApplicationCtxKey).(*appDomain.Application)
		if !ok {
			handleError(w, fmt.Errorf("missing application"))
			return
		}

		if err := json.NewEncoder(w).Encode(app); err != nil {
			handleError(w, err)
			return
		}
	}
}

func showSelectFilters(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func saveAppFilters(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")

		app, ok := r.Context().Value(ApplicationCtxKey).(*domain.Application)
		if !ok {
			handleError(w, fmt.Errorf("unknown application"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(app); err != nil {
			log.Printf("receiving new app filters: %v\n", err)
			handleError(w, err)
			return
		}

		if err := svc.SaveApplicationFilters(r.Context(), app); err != nil {
			log.Printf("saving filters: %v\n", err)
			handleError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(app); err != nil {
			handleError(w, err)
			return
		}
	}
}

func saveApp(svc application.Service, iFactsCli iFacts.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")

		var p appDomain.Application
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			log.Printf("receiving new app data: %v\n", err)
			handleError(w, err)
			return
		}

		fmt.Printf("Request to save application: %#v\n", p)

		if saved, err := svc.SaveApplicationOrImportFromIFacts(r.Context(), &p, iFactsCli); err != nil {
			handleError(w, err)
			return
		} else {
			p = *saved
		}

		fmt.Printf("application saved or found: %#v\n", p)

		w.WriteHeader(http.StatusCreated)
		fmt.Printf("Respond with saved application: %#v\n", p)
		if err := json.NewEncoder(w).Encode(p); err != nil {
			handleError(w, err)
			return
		}
	}
}

func updateApp(svc application.Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		app, ok := r.Context().Value(ApplicationCtxKey).(*appDomain.Application)
		if !ok {
			handleError(w, fmt.Errorf("missing application"))
			return
		}

		w.Header().Set("Content-type", "application/json")

		var p appDomain.Application
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			handleError(w, err)
			return
		}

		p.ID = app.ID

		if err := svc.SaveApplication(r.Context(), &p); err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func controlsForApp(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		app, ok := r.Context().Value(ApplicationCtxKey).(*appDomain.Application)
		if !ok {
			handleError(w, fmt.Errorf("missing application"))
			return
		}

		ctrls, err := svc.FetchControlsByApplication(r.Context(), app)
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ctrls); err != nil {
			log.Printf("Encoding controls for app: %v\n", err)
			handleError(w, err)
			return
		}
	}
}

func previewControlsForApp(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app, ok := r.Context().Value(ApplicationCtxKey).(*appDomain.Application)
		if !ok {
			handleError(w, fmt.Errorf("missing application"))
			return
		}

		ctrls, err := svc.FilterControls(r.Context(), appToFilter(app))
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ctrls); err != nil {
			log.Printf("Encoding controls for app: %v\n", err)
			handleError(w, err)
			return
		}
	}
}

func saveControlsForApp(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		app, ok := r.Context().Value(ApplicationCtxKey).(*appDomain.Application)
		if !ok {
			handleError(w, fmt.Errorf("missing application"))
			return
		}

		controls, err := svc.FilterControls(r.Context(), appToFilter(app))
		if err != nil {
			handleError(w, err)
			return
		}

		err = svc.SaveControlsForApplication(r.Context(), app, controls)
		if err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func showAppControlDetails(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app, ok := r.Context().Value(ApplicationCtxKey).(*appDomain.Application)
		if !ok {
			handleError(w, fmt.Errorf("missing application"))
			return
		}

		ctrlID := chi.URLParam(r, "id")

		ctrl, err := svc.FetchAppControlByID(r.Context(), app, ctrlID)
		if err != nil {
			handleError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(ctrl); err != nil {
			log.Printf("Encoding control for app: %v\n", err)
			handleError(w, err)
			return
		}
	}
}
func appToFilter(app *domain.Application) domain.ControlsFilter {
	filter := domain.ControlsFilter{}

	if app.OnlyHandledCentrally {
		filter.OnlyHandleCentrally = &app.OnlyHandledCentrally
	}
	if app.HandledCentrallyBy != "" {
		filter.HandledCentrallyBy = &app.HandledCentrallyBy
	}
	if app.ExcludeForExternalSupplier {
		filter.ExcludeForExternalSupplier = &app.ExcludeForExternalSupplier
	}
	if app.SoftwareDevelopmentRelevant {
		filter.SoftwareDevelopmentRelevant = &app.SoftwareDevelopmentRelevant
	}
	if app.CloudOnly {
		filter.CloudOnly = &app.CloudOnly
	}
	if app.PhysicalSecurityOnly {
		filter.PhysicalSecurityOnly = &app.PhysicalSecurityOnly
	}
	if app.PersonalSecurityOnly {
		filter.PersonalSecurityOnly = &app.PersonalSecurityOnly
	}

	return filter
}

func saveLocallyFromIFacts(svc application.Service, ifclient iFacts.Client) http.HandlerFunc {

	type request struct {
		InternalID int64  `json:"app_internal_id"`
		Name       string `json:"app_name"`
		IFactsID   string `json:"ifacts_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		iFactsID := chi.URLParam(r, "iFactsID")
		if iFactsID == "" {
			handleError(w, fmt.Errorf("missing iFacts ID"))
			return
		}

		var b request
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			handleError(w, err)
			return
		}

		classifications, err := ifclient.GetClassifications(iFactsID)
		if err != nil {
			handleError(w, err)
			return
		}

		svc.SaveFromIFacts(r.Context(), iFactsID, ifclient)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(classifications); err != nil {
			handleError(w, err)
		}
	}
}
