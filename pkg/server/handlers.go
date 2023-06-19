package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mehix/sec-checklist/pkg/iFacts"
)

func configIFactsClient(ifc *iFacts.Client) http.HandlerFunc {
	type req struct {
		CookieToken string `json:"cookie_token"`
		BaseURL     string `json:"base_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var payload req
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			handleError(w, err)
			return
		}
		ifc.BaseURL = payload.BaseURL
		ifc.CookieToken = payload.CookieToken

		w.WriteHeader(http.StatusCreated)
	}
}

func forwardGetToIFacts(ifc *iFacts.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		copyResponse := func(resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("iFacts response: %s", resp.Status)
			}

			_, err := io.Copy(w, resp.Body)
			return err
		}

		if err := ifc.Request(http.MethodGet, r.URL.Path, nil, copyResponse); err != nil {
			handleError(w, err)
		}
	}
}

func searchAppByNameRemote(ifc *iFacts.Client) http.HandlerFunc {
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
