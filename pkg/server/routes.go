package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mehix/sec-checklist/pkg/domain/check"
	"github.com/mehix/sec-checklist/pkg/iFacts"
	"github.com/mehix/sec-checklist/pkg/service/checks"
)

func Handlers(svc checks.Service, iFactsClient *iFacts.Client) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.AllowContentEncoding("application/json"))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//r.Use(middleware.RedirectSlashes)

	//r.Get("/apps/", listAllApps(svc))
	//r.Post("/apps", saveApp(svc))
	//r.Get("/apps/{id:[0-9]+}", showAppByID(svc))

	r.Get("/controls/", showAll(svc))
	r.Post("/controls/", showFiltered(svc))
	r.Get("/controls/{id:[0-9.]+}", showOne(svc))

	r.Post("/ifacts/config", configIFactsClient(iFactsClient))
	r.Method(http.MethodGet, "/ifacts/*", http.StripPrefix("/ifacts", http.HandlerFunc(forwardGetToIFacts(iFactsClient))))

	r.Get("/docs/controls/filter", showFiltered(svc))

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

	exampleFilter := func(w http.ResponseWriter) {
		t := true
		s := "BSO"
		json.NewEncoder(w).Encode(struct {
			Req  filter          `json:"Request body example"`
			Resp []check.Control `json:"Response example"`
		}{Req: filter{
			OnlyHandleCentrally:         &t,
			HandledCentrallyBy:          &s,
			ExcludeForExternalSupplier:  &t,
			SoftwareDevelopmentRelevant: &t,
			CloudOnly:                   &t,
			PhysicalSecurityOnly:        &t,
			PersonalSecurityOnly:        &t,
		},
			Resp: []check.Control{{}}})
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			// show an example of filter
			exampleFilter(w)
			return
		}

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
		var filtered []check.Control
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
