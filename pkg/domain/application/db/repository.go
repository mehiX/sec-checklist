package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mehix/sec-checklist/pkg/domain"
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

func (r *repository) FetchByID(ctx context.Context, id string) (*domain.Application, error) {
	sql := "select * from V_APPS where ID=?"
	row := r.db.QueryRowContext(ctx, sql, id)

	return scanForApp(row)
}

func (r *repository) FindByInternalID(ctx context.Context, internalID int) (*domain.Application, error) {

	sql := "select * from V_APPS where internal_id=?"
	row := r.db.QueryRowContext(ctx, sql, internalID)

	return scanForApp(row)
}

func (r *repository) ListAll(ctx context.Context) ([]domain.Application, error) {
	qry := "select * from V_APPS"

	rows, err := r.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []domain.Application{}
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

func (r *repository) SaveFilters(ctx context.Context, app *domain.Application) error {

	updateAppProfile := `update APP_PROFILES set
	 only_handle_centrally = ?,
	 handled_centrally_by = ?,
	 excluded_for_external_supplier = ?,
	 software_development_relevant = ?,
     cloud_only = ?,
	 physical_security_only = ?,
	 personal_security_only = ? 
	 where APP_ID = ?`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx, updateAppProfile,
		app.OnlyHandledCentrally,
		app.HandledCentrallyBy, app.ExcludeForExternalSupplier,
		app.SoftwareDevelopmentRelevant, app.CloudOnly,
		app.PhysicalSecurityOnly, app.PersonalSecurityOnly, app.ID)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			log.Printf("failed rollback for tx when updating app_profiles: %v\n", err)
		}
		return err
	}

	return tx.Commit()

}

func (r *repository) Save(ctx context.Context, app *domain.Application) (err error) {
	insertApp := "insert into APPS (id, internal_id, name, ifacts_id) values (?, ?, ?, ?)"
	insertEmptyAppFilters := "insert into APPS_PROFILES (APP_ID) values (?)"

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

	_, err = tx.ExecContext(ctx, insertApp, app.ID, app.InternalID, app.Name, app.IFactsID)
	if err != nil {
		return
	}

	_, err = tx.ExecContext(ctx, insertEmptyAppFilters, app.ID)
	if err != nil {
		log.Printf("failed to insert empty app_profiles records: %v\n", err)
	}

	return nil
}

func (r *repository) Update(ctx context.Context, app *domain.Application) (err error) {
	updateApp := "update APPS set name = ? where id = ?)"

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

	return nil
}

func (r *repository) SaveIFactsClassifications(ctx context.Context, id string, classifications []domain.Classification) error {

	sqlInsert := `insert into APP_CLASSIFICATIONS 
	(IFACTS_ID, classification_id, classification_name,
		level_id, level_name) values (?, ?, ?, ?, ?)`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	prep, err := tx.PrepareContext(ctx, sqlInsert)
	if err != nil {
		return err
	}
	for _, clsf := range classifications {
		if _, err := prep.ExecContext(ctx, id, clsf.ID, clsf.Name, clsf.LevelID, clsf.LevelName); err != nil {
			log.Printf("error saving classification: %v\n", err)
		}
	}

	return tx.Commit()
}
