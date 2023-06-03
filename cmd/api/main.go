package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mehix/sec-checklist/pkg/server"
	"github.com/mehix/sec-checklist/pkg/service/checks"
	"golang.org/x/net/context"
)

var fromExcel = flag.String("from", "", "Path to Excel file to read initial data")
var fromSheet = flag.String("fromSheet", "", "Name of the Excel sheet (default is the first sheet)")
var addr = flag.String("http", "", "HTTP address to listen on")
var db = flag.Bool("db", false, "Connect to database")
var initDB = flag.Bool("init", false, "Initialize the database with data from the Excel file")

func init() {
	godotenv.Load()
	flag.Parse()
}

func main() {
	svc := checks.NewService()

	if *fromExcel != "" {
		checks.WithXls(*fromExcel, *fromSheet)(svc)
	}

	if *db {
		checks.WithDb(os.Getenv("CHECKLISTS_DSN"))(svc)
	}

	if *initDB {
		fmt.Printf("Read Excel file: %s\n", *fromExcel)
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

	if *addr != "" {
		fmt.Println("Listening on", *addr)
		if err := http.ListenAndServe(*addr, server.Handlers(svc)); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
