package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mehix/sec-checklist/pkg/iFacts"
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
