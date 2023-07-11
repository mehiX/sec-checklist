package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	appDomain "github.com/mehix/sec-checklist/pkg/domain/application"
	"github.com/mehix/sec-checklist/pkg/service/application"
)

func listAllApps(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		apps, err := svc.ListAll(r.Context())
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

func saveApp(svc application.Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")

		var p appDomain.Application
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			log.Printf("receiving new app data: %v\n", err)
			handleError(w, err)
			return
		}

		fmt.Printf("Request to save application: %#v\n", p)

		if err := svc.Save(r.Context(), &p); err != nil {
			handleError(w, err)
			return
		}

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

		if err := svc.Update(r.Context(), &p); err != nil {
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

		ctrls, err := svc.FetchByApplication(r.Context(), app)
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
