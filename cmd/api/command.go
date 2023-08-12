package api

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/mehix/sec-checklist/pkg/iFacts"
	"github.com/mehix/sec-checklist/pkg/server"
	"github.com/mehix/sec-checklist/pkg/service/application"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	addr string
	noDb bool
)

var fromExcel, fromSheet string

func init() {
	godotenv.Load()
	flag.Parse()
}

func Command() *cobra.Command {
	cmdApi := &cobra.Command{
		Use:   "api",
		Short: "Start or manage the API",
	}

	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Starts the backend REST API",
		Long:  "Runs an HTTP server and exposes a REST API",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}

	cmdServe.Flags().StringVar(&addr, "http", "127.0.0.1:8080", "HTTP address to listen on")
	cmdServe.Flags().BoolVar(&noDb, "no-db", false, "do not try to connect to DB")

	cmdLoad := &cobra.Command{
		Use:   "load",
		Short: "Loads data into the database",
		Long:  "Reads data from an Excel file and loads it into the database",
		Run: func(cmd *cobra.Command, args []string) {
			importData()
		},
	}

	cmdLoad.Flags().StringVar(&fromExcel, "from", "", "Path to Excel file to read initial data")
	cmdLoad.Flags().StringVar(&fromSheet, "fromSheet", "", "Name of the Excel sheet (default is the first sheet)")

	cmdApi.AddCommand(cmdServe, cmdLoad)

	return cmdApi
}

func serve() {
	var wg sync.WaitGroup
	wg.Add(2)

	svcApps := application.NewService()

	go func() {
		defer wg.Done()

		application.WithControlsDb(os.Getenv("CHECKLISTS_DSN"))(svcApps)
	}()

	go func() {
		defer wg.Done()

		application.WithAppsDb(os.Getenv("CHECKLISTS_DSN"))(svcApps)
	}()

	wg.Wait()

	iFactsClient := iFacts.NewClient(
		os.Getenv("IFACTS_BASEURL"),
		os.Getenv("IFACTS_CLIENT_ID"),
		os.Getenv("IFACTS_CLIENT_SECRET"),
	)

	fmt.Println("Listening on", addr)
	h := server.Handlers(svcApps, iFactsClient)
	printRoutesHelp(h)
	if err := http.ListenAndServe(addr, h); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func importData() {
	svc := application.NewService(
		application.WithControlsDb(os.Getenv("CHECKLISTS_DSN")),
		application.WithXls(fromExcel, fromSheet),
	)

	fmt.Printf("Read Excel file: %s\n", fromExcel)
	ctrls, err := svc.FetchControlsFromExcel()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d controls and categories\n", len(ctrls))
	if len(ctrls) > 0 {
		fmt.Println(ctrls[0])
	}

	if err := svc.SaveAllControls(context.Background(), ctrls); err != nil {
		log.Printf("Saving checks to database: %v\n", err)
	}

}

func printRoutesHelp(h http.Handler) {
	r, ok := h.(*chi.Mux)
	if !ok {
		return
	}
	var out = tabwriter.NewWriter(os.Stdout, 10, 8, 0, '\t', 0)
	printHelp(out, "", r.Routes())
	out.Flush()

}

func printHelp(out *tabwriter.Writer, parentPattern string, routes []chi.Route) {

	fmt.Fprintln(out)

	for _, r := range routes {
		ptrn := strings.TrimSuffix(r.Pattern, "/*")
		if r.SubRoutes != nil {
			printHelp(out, parentPattern+ptrn, r.SubRoutes.Routes())
		} else {
			for m := range r.Handlers {
				fmt.Fprintf(out, "[%s]\t%s\n", m, parentPattern+ptrn)
			}
			fmt.Fprintln(out)
		}

	}

}
