package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mehix/sec-checklist/pkg/iFacts"
	"github.com/mehix/sec-checklist/pkg/service/application"
)

func searchIFactsAppByName(ifc iFacts.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		if q == "" {
			w.Header().Set("Content-type", "application/json")
			w.Write([]byte(`{"Assets": []}`))
			return
		}

		cpResp := func(resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("iFacts response: %s", resp.Status)
			}

			_, err := io.Copy(w, resp.Body)
			return err
		}

		w.Header().Set("Content-Type", "application/json-patch+json")
		if err := ifc.SearchByName(q, cpResp); err != nil {
			handleError(w, err)
		}
	}
}

func getIFactsClassifications(svc application.Service, ifc iFacts.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		iFactsID := chi.URLParam(r, "iFactsID")

		classifications, err := ifc.GetClassifications(iFactsID)
		if err != nil {
			handleError(w, err)
			return
		}

		if err := svc.SaveFromIFacts(r.Context(), iFactsID, ifc); err != nil {
			log.Printf("saving data from iFacts: %v\n", err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(classifications); err != nil {
			handleError(w, err)
		}

	}
}
