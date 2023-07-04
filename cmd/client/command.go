package client

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var addr, apiAddr string

//go:embed templates/*.tmpl
var fs embed.FS

type appIDCtxKey struct{}

var AppIDCtxKey = &appIDCtxKey{}

func Command() *cobra.Command {

	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "User interface",
		Long:  "Starts a frontend application to interact with users",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}

	clientCmd.Flags().StringVar(&addr, "http", "127.0.0.1:8081", "http port to listen for requests")
	clientCmd.Flags().StringVar(&apiAddr, "api", "http://127.0.0.1:8080", "URL to the backend REST API")

	return clientCmd
}

func serve() {
	tmpl, err := template.ParseFS(fs, "templates/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/apps")
		w.WriteHeader(http.StatusFound)
	})
	http.HandleFunc("/apps/new", showTemplate(tmpl, "filters"))
	http.HandleFunc("/apps/filters", showTemplate(tmpl, "filters"))
	http.HandleFunc("/apps/", showApps(tmpl))

	http.HandleFunc("/config", showTemplate(tmpl, "config"))

	fmt.Printf("Listening on %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func showApps(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/apps/")
		if id != "" {
			ctx := context.WithValue(r.Context(), AppIDCtxKey, id)
			r = r.WithContext(ctx)
		}

		showTemplate(tmpl, "apps")(w, r)
	}
}

func showTemplate(t *template.Template, name string) http.HandlerFunc {
	type data struct {
		ApiURL        string
		SelectedAppID string
	}

	d := data{
		ApiURL: apiAddr,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html; charset=utf-8")

		if appID, ok := r.Context().Value(AppIDCtxKey).(string); ok {
			log.Println("App ID", appID)
			d.SelectedAppID = appID
		}

		if err := t.ExecuteTemplate(w, name, d); err != nil {
			log.Printf("Executing template index: %v\n", err.Error())
			w.Write([]byte(err.Error()))
			return
		}
	}
}
