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
		log.Println("APPS DB connection failed", err)
	} else {
		fmt.Println("Connected to APPS DB")
	}

	return &repository{db: db}
}

func (r *repository) FetchByID(ctx context.Context, id string) (*application.Application, error) {
	sql := "select * from V_APPS where ID=?"
	row := r.db.QueryRowContext(ctx, sql, id)

	return scanForApp(row)
}

func (r *repository) ListAll(ctx context.Context) ([]application.Application, error) {
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

func (r *repository) Save(ctx context.Context, app *application.Application) (err error) {
	insertApp := "insert into APPS (id, internal_id, name) values (?, ?, ?)"
	insertAppProfile := `insert into APP_PROFILES (
		APP_ID, only_handle_centrally,handled_centrally_by,excluded_for_external_supplier,software_development_relevant,
        cloud_only,physical_security_only,personal_security_only) values (?, ?, ?, ?, ?, ?, ?, ?)`
	insertAppClassifications := `insert into APP_CLASSIFICATIONS (
		APP_ID, c, i, a, t) values (?, ?, ?, ?, ?)`

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

	_, err = tx.ExecContext(ctx, insertApp, app.ID, app.InternalID, app.Name)
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

	_, err = tx.ExecContext(
		ctx, insertAppClassifications,
		app.ID, app.C, app.I, app.A, app.T)
	if err != nil {
		return
	}

	return nil
}

func (r *repository) Update(ctx context.Context, app *application.Application) (err error) {
	updateApp := "update APPS set name = ? where id = ?)"
	updateAppProfile := `update APP_PROFILES set
	 only_handle_centrally = ?,
	 handled_centrally_by = ?,
	 excluded_for_external_supplier = ?,
	 software_development_relevant = ?,
     cloud_only = ?,
	 physical_security_only = ?,
	 personal_security_only = ? 
	 where APP_ID = ?`
	updateAppClassifications := `update APP_CLASSIFICATIONS set c = ?, i = ?, a = ?, t = ? where APP_ID = ?`

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

	_, err = tx.ExecContext(ctx, updateApp, app.Name, app.ID)
	if err != nil {
		return
	}
	_, err = tx.ExecContext(
		ctx, updateAppProfile,
		app.OnlyHandledCentrally,
		app.HandledCentrallyBy, app.ExcludeForExternalSupplier,
		app.SoftwareDevelopmentRelevant, app.CloudOnly,
		app.PhysicalSecurityOnly, app.PersonalSecurityOnly,
		app.ID)
	if err != nil {
		return
	}

	_, err = tx.ExecContext(
		ctx, updateAppClassifications,
		app.C, app.I, app.A, app.T, app.ID)
	if err != nil {
		return
	}

	return nil
}
