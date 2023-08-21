package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mehix/sec-checklist/pkg/domain"
	"github.com/mehix/sec-checklist/pkg/iFacts"
	"github.com/mehix/sec-checklist/pkg/service/application"
)

func Handlers(svcApps application.Service, iFactsClient iFacts.Client) http.Handler {
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

	r.Route("/apps", func(r chi.Router) {
		r.Get("/", listAllApps(svcApps))
		r.Post("/", saveApp(svcApps, iFactsClient))
		r.Route("/{id:[0-9a-zA-Z-]+}", func(r chi.Router) {
			r.Use(ApplicationCtx(svcApps))
			r.Get("/", showAppByID(svcApps))
			r.Put("/", updateApp(svcApps))
			r.Get("/filters", showSelectFilters(svcApps)) // TODO
			r.Post("/filters", saveAppFilters(svcApps))
			r.Get("/controls", controlsForApp(svcApps))
			r.Get("/controls/preview", previewControlsForApp(svcApps))
		})
		r.Get("/iFacts/search", searchIFactsAppByName(iFactsClient))
		r.Get("/iFacts/classifications/{iFactsID:[0-9a-zA-Z-]+}", getIFactsClassifications(svcApps, iFactsClient))
	})

	r.Route("/controls", func(r chi.Router) {
		r.Get("/", showAll(svcApps))
		r.Post("/filter/", showFiltered(svcApps))
		r.Get("/{id:[0-9.]+}", showOneControl(svcApps))
	})

	//r.Post("/ifacts/config", configIFactsClient(iFactsClient))
	// forward a request to iFacts and return the result
	//r.Method(http.MethodGet, "/ifacts/*", http.StripPrefix("/ifacts", http.HandlerFunc(forwardGetToIFacts(iFactsClient))))

	r.Get("/docs/controls/filter", showFiltered(svcApps))

	return r
}

type applicationCtxKey struct{}

var ApplicationCtxKey = &applicationCtxKey{}

func ApplicationCtx(svc application.Service) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appID := chi.URLParam(r, "id")
			app, err := svc.FetchApplicationByID(r.Context(), appID)
			if err != nil {
				handleError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), ApplicationCtxKey, app)
			h.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func showOneControl(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParam(r, "id")
		one, err := svc.FetchControlByID(r.Context(), id)
		if err != nil {
			handleError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(one); err != nil {
			handleError(w, err)
		}
	}
}
func showAll(svc application.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		all, err := svc.FetchAllControls()
		if err != nil {
			handleError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(all); err != nil {
			handleError(w, err)
		}
	}
}

func showFiltered(svc application.Service) http.HandlerFunc {

	exampleFilter := func(w http.ResponseWriter) {
		t := true
		s := "BSO"
		json.NewEncoder(w).Encode(struct {
			Req  domain.ControlsFilter `json:"Request body example"`
			Resp []domain.Control      `json:"Response example"`
		}{Req: domain.ControlsFilter{
			OnlyHandleCentrally:         &t,
			HandledCentrallyBy:          &s,
			ExcludeForExternalSupplier:  &t,
			SoftwareDevelopmentRelevant: &t,
			CloudOnly:                   &t,
			PhysicalSecurityOnly:        &t,
			PersonalSecurityOnly:        &t,
		},
			Resp: []domain.Control{{}}})
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			// show an example of filter
			exampleFilter(w)
			return
		}

		var payload domain.ControlsFilter
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			handleError(w, err)
			return
		}

		// TODO: move the filter to the service
		all, err := svc.FetchAllControls()
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
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("handling error: %v\n", err)
	}
}
