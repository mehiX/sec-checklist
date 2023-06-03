package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mehix/sec-checklist/pkg/domain"
	"github.com/mehix/sec-checklist/pkg/service/checks"
)

func Handlers(svc checks.Service) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.AllowContentEncoding("application/json"))
	r.Use(middleware.RedirectSlashes)

	r.Get("/controls", showAll(svc))
	r.Post("/controls", showFiltered(svc))
	r.Get("/controls/{id:[0-9.]+}", showOne(svc))

	return r
}

func showOne(svc checks.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParam(r, "id")
		one, err := svc.FetchByID(r.Context(), id)
		if err != nil {
			handleError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(one); err != nil {
			handleError(w, err)
		}
	}
}
func showAll(svc checks.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		all, err := svc.FetchAll()
		if err != nil {
			handleError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(all); err != nil {
			handleError(w, err)
		}
	}
}

func showFiltered(svc checks.Service) http.HandlerFunc {

	type filter struct {
		OnlyHandleCentrally         *bool   `json:"only_handle_centrally,omitempty"`
		HandledCentrallyBy          *string `json:"handled_centrally_by,omitempty"`
		ExcludeForExternalSupplier  *bool   `json:"exclude_for_external_supplier,omitempty"`
		SoftwareDevelopmentRelevant *bool   `json:"software_development_relevant,omitempty"`
		CloudOnly                   *bool   `json:"cloud_only,omitempty"`
		PhysicalSecurityOnly        *bool   `json:"physical_security_only,omitempty"`
		PersonalSecurityOnly        *bool   `json:"personal_security_only,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var payload filter
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			handleError(w, err)
			return
		}

		// TODO: move the filter to the service
		all, err := svc.FetchAll()
		if err != nil {
			handleError(w, err)
			return
		}
		var filtered []domain.Control
		for _, c := range all {
			if payload.OnlyHandleCentrally != nil && *payload.OnlyHandleCentrally != c.OnlyHandledCentrally {
				continue
			}
			if payload.HandledCentrallyBy != nil && !strings.Contains(c.HandledCentrallyBy, *payload.HandledCentrallyBy) {
				continue
			}
			if payload.ExcludeForExternalSupplier != nil && *payload.ExcludeForExternalSupplier != c.ExcludeForExternalSupplier {
				continue
			}
			if payload.SoftwareDevelopmentRelevant != nil && *payload.SoftwareDevelopmentRelevant != c.SoftwareDevelopmentRelevant {
				continue
			}
			if payload.CloudOnly != nil && *payload.CloudOnly != c.CloudOnly {
				continue
			}
			if payload.PhysicalSecurityOnly != nil && *payload.PhysicalSecurityOnly != c.PhysicalSecurityOnly {
				continue
			}
			if payload.PersonalSecurityOnly != nil && *payload.PersonalSecurityOnly != c.PersonalSecurityOnly {
				continue
			}

			filtered = append(filtered, c)
		}

		if err := json.NewEncoder(w).Encode(filtered); err != nil {
			handleError(w, err)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {

	resp := struct {
		Error string
	}{
		Error: err.Error(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("handling error: %v\n", err)
	}
}
