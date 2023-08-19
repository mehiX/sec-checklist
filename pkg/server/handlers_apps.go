package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
