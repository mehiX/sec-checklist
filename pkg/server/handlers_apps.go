package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

		appID := chi.URLParam(r, "id")
		app, err := svc.FetchByID(r.Context(), appID)
		if err != nil {
			handleError(w, err)
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
			handleError(w, err)
			return
		}

		if err := svc.Save(r.Context(), &p); err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func updateApp(svc application.Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")

		var p appDomain.Application
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			handleError(w, err)
			return
		}

		if err := svc.Update(r.Context(), &p); err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
