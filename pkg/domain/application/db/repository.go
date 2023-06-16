package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mehix/sec-checklist/pkg/domain/application"
	"github.com/mehix/sec-checklist/pkg/domain/db"
)

type repository struct {
	db *sql.DB
}

func NewRepository(dsn string) application.ReaderWriter {
	db, err := db.ConnWithRetry(db.Conn, 5, time.Second, time.Minute)(context.Background(), dsn)
	if err != nil {
		log.Fatalln("DB connection failed", err)
	}

	fmt.Println("Connected to database")
	return &repository{db: db}
}

func (r *repository) FetchAppByID(ctx context.Context, id string) (*application.Application, error) {
	sql := "select * from V_APPS where ID=?"
	row := r.db.QueryRowContext(ctx, sql, id)

	return scanForApp(row)
}

func (r *repository) ListAllApps(ctx context.Context) ([]application.Application, error) {
	qry := "select * from V_APPS"

	rows, err := r.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []application.Application{}
	for rows.Next() {
		app, err := scanForApp(rows)
		if err != nil {
			log.Printf("Scanning row for ListAllApps: %v\n", err)
			continue
		}
		apps = append(apps, *app)
	}

	return apps, nil
}

func (r *repository) SaveApp(ctx context.Context, app *application.Application) (err error) {
	insertApp := "insert into APPS (id, name) values (?, ?)"
	insertAppProfile := `insert into APP_PROFILES (
		APP_ID, only_handle_centrally,handled_centrally_by,excluded_for_external_supplier,software_development_relevant,
        cloud_only,physical_security_only,personal_security_only) values (?, ?, ?, ?, ?, ?, ?, ?)`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("Rolling back TX: %v\n", err.Error())
			}
			return
		}
		if err := tx.Commit(); err != nil {
			log.Printf("Committing TX: %v\n", err.Error())
		}
	}()

	_, err = tx.ExecContext(ctx, insertApp, app.ID, app.Name)
	if err != nil {
		return
	}
	_, err = tx.ExecContext(
		ctx, insertAppProfile,
		app.ID, app.OnlyHandledCentrally,
		app.HandledCentrallyBy, app.ExcludeForExternalSupplier,
		app.SoftwareDevelopmentRelevant, app.CloudOnly,
		app.PhysicalSecurityOnly, app.PersonalSecurityOnly)
	if err != nil {
		return
	}

	return nil
}
