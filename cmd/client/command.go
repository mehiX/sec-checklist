package client

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/spf13/cobra"
)

var addr, apiAddr string

//go:embed templates/*.tmpl
var fs embed.FS

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
	http.HandleFunc("/apps/", showTemplate(tmpl, "apps"))
	http.HandleFunc("/apps/new", showTemplate(tmpl, "filters"))
	http.HandleFunc("/apps/filters", showTemplate(tmpl, "filters"))

	http.HandleFunc("/config", showTemplate(tmpl, "config"))

	fmt.Printf("Listening on %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func showTemplate(t *template.Template, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		if err := t.ExecuteTemplate(w, name, nil); err != nil {
			log.Printf("Executing template index: %v\n", err.Error())
			w.Write([]byte(err.Error()))
			return
		}
	}
}
