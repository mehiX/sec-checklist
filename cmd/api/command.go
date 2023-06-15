package api

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mehix/sec-checklist/pkg/iFacts"
	"github.com/mehix/sec-checklist/pkg/server"
	"github.com/mehix/sec-checklist/pkg/service/checks"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var addr string
var db bool

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
	cmdServe.Flags().BoolVar(&db, "db", false, "Connect to database")

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
	svc := checks.NewService()
	iFactsClient := &iFacts.Client{}

	if db {
		checks.WithDb(os.Getenv("CHECKLISTS_DSN"))(svc)
	}

	fmt.Println("Listening on", addr)
	if err := http.ListenAndServe(addr, server.Handlers(svc, iFactsClient)); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func importData() {
	svc := checks.NewService(
		checks.WithDb(os.Getenv("CHECKLISTS_DSN")),
		checks.WithXls(fromExcel, fromSheet),
	)

	fmt.Printf("Read Excel file: %s\n", fromExcel)
	ctrls, err := svc.FetchAllFromExcel()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d controls and categories\n", len(ctrls))
	if len(ctrls) > 0 {
		fmt.Println(ctrls[0])
	}

	if err := svc.SaveAll(context.Background(), ctrls); err != nil {
		log.Printf("Saving checks to database: %v\n", err)
	}

}
