package client

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
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

	fmt.Printf("Listening on %s\n", addr)

	if err := http.ListenAndServe(addr, Handlers(tmpl)); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func Handlers(tmpl *template.Template) http.Handler {
	r := chi.NewMux()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/apps")
		w.WriteHeader(http.StatusFound)
	})
	r.Route("/apps", func(r chi.Router) {
		r.Get("/new", showTemplate(tmpl, "step1"))
		r.Get("/", showApps(tmpl))
		r.Route("/{id:[0-9a-zA-Z-]+}", func(r chi.Router) {
			r.Use(AppCtx)
			r.Get("/", showApps(tmpl))
			r.Get("/filters", showTemplate(tmpl, "step2"))
			r.Get("/controls", showTemplate(tmpl, "step3"))
			r.Get("/controls/{id:[0-9.]+}/", showApplicationControl(tmpl))
		})
	})

	return r
}

func AppCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id != "" {
			ctx := context.WithValue(r.Context(), AppIDCtxKey, id)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func showApps(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func showApplicationControl(t *template.Template) http.HandlerFunc {
	type data struct {
		ApiURL        string
		SelectedAppID string
		ControlID     string
	}

	d := data{
		ApiURL: apiAddr,
	}

	templateName := "appControl"

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html; charset=utf-8")

		if appID, ok := r.Context().Value(AppIDCtxKey).(string); ok {
			d.SelectedAppID = appID
		}

		ctrlID := chi.URLParam(r, "id")
		if ctrlID == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("missing control ID"))
			return
		}

		d.ControlID = ctrlID

		if err := t.ExecuteTemplate(w, templateName, d); err != nil {
			log.Printf("Executing template index: %v\n", err.Error())
			w.Write([]byte(err.Error()))
			return
		}

	}
}
