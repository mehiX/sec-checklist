package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mehix/sec-checklist/pkg/domain/application"
	"github.com/mehix/sec-checklist/pkg/domain/check"
	"github.com/mehix/sec-checklist/pkg/domain/db"
)

type repository struct {
	db *sql.DB
}

func NewRepository(dsn string) check.ReaderWriter {
	db, err := db.ConnWithRetry(db.Conn, 5, time.Second, time.Minute)(context.Background(), dsn)
	if err != nil {
		log.Println("CHECKS DB connection failed", err)
	} else {
		fmt.Println("Connected to CHECKS DB")
	}

	return &repository{db: db}
}

func (r *repository) FetchAll() ([]check.Control, error) {
	qry := "select * from V_CHECKS"

	rows, err := r.db.QueryContext(context.TODO(), qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ctrls []check.Control
	for rows.Next() {
		ctrl, err := scanForControl(rows)
		if err != nil {
			log.Printf("Scanning SQL row: %v\n", err)
			continue
		}
		ctrls = append(ctrls, ctrl)
	}

	return ctrls, nil

}
func (r *repository) FetchByType(t string) ([]check.Control, error) { return nil, nil }

func (r *repository) FetchByID(ctx context.Context, id string) (check.Control, error) {
	qry := "select * from V_CHECKS where ID = ?"

	row := r.db.QueryRowContext(ctx, qry, id)
	return scanForControl(row)
}

func (r *repository) SaveAll(ctx context.Context, all []check.Control) (err error) {
	qry := "insert into CHECKS (ID, type, name, description,asset_type,last_update,old_id) values (?, ?, ?, ?, ?, ?, ?)"
	qryCiat := "insert CHECKS_CIAT (CHECK_ID, c, i, a, t) values (?, ?, ?, ?, ?)"
	qryFilters := `insert FILTERS (CHECK_ID, only_handle_centrally, 
		handled_centrally_by, excluded_for_external_supplier, 
		software_development_relevant, cloud_only, 
		physical_security_only, personal_security_only) 
		values (?, ?, ?, ?, ?, ?, ?, ?)`
	qryOthers := `insert CHECKS_OTHERS (CHECK_ID, pd, nsi, sese, otcl, 
		csr, spsa, spsa_unique, gdpr, gdpr_unique, external_supplier,
		operational_capability, part_of_gisr) 
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var tx *sql.Tx

	tx, err = r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("Rollback TX for SaveAll: %v\n", err)
				return
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Printf("Commit TX for SaveAll: %v\n", err)
			}
		}
	}()

	var stmt *sql.Stmt
	stmt, err = tx.PrepareContext(context.TODO(), qry)
	if err != nil {
		return
	}
	defer stmt.Close()

	var stmtCiat *sql.Stmt
	stmtCiat, err = tx.PrepareContext(context.TODO(), qryCiat)
	if err != nil {
		return
	}
	defer stmtCiat.Close()

	var stmtFilters *sql.Stmt
	stmtFilters, err = tx.PrepareContext(context.TODO(), qryFilters)
	if err != nil {
		return
	}

	var stmtOthers *sql.Stmt
	stmtOthers, err = tx.PrepareContext(context.TODO(), qryOthers)
	if err != nil {
		return
	}

	for _, c := range all {
		_, err = stmt.ExecContext(context.TODO(), c.ID, c.Type, c.Name, c.Description, c.AssetType, c.LastUpdated, c.OldID)
		if err != nil {
			return
		}

		_, err = stmtCiat.ExecContext(context.TODO(), c.ID, c.C, c.I, c.A, c.T)
		if err != nil {
			return
		}

		_, err = stmtFilters.ExecContext(
			context.TODO(), c.ID, c.OnlyHandledCentrally,
			c.HandledCentrallyBy, c.ExcludeForExternalSupplier,
			c.SoftwareDevelopmentRelevant, c.CloudOnly,
			c.PhysicalSecurityOnly, c.PersonalSecurityOnly)
		if err != nil {
			return
		}

		_, err = stmtOthers.ExecContext(
			context.TODO(), c.ID, c.PD, c.NSI, c.SESE, c.OTCL, c.CSRDirection,
			c.SPSA, c.SPSAUnique, c.GDPR, c.GDPRUnique, c.ExternalSupplier,
			c.OperationalCapability, c.PartOfGISR)
		if err != nil {
			return
		}

	}

	return nil
}

func (r *repository) FetchForApplication(ctx context.Context, app *application.Application) ([]check.Control, error) {
	qry := "select * from V_CHECKS where 1=1"
	args := make([]any, 0)

	if app.OnlyHandledCentrally {
		qry += " AND only_handle_centrally = ?"
		args = append(args, app.OnlyHandledCentrally)
	}
	if app.HandledCentrallyBy != "" {
		qry += " AND handled_centrally_by = ?"
		args = append(args, app.HandledCentrallyBy)
	}
	if app.ExcludeForExternalSupplier {
		qry += " AND excluded_for_external_supplier = ?"
		args = append(args, app.ExcludeForExternalSupplier)
	}
	if app.SoftwareDevelopmentRelevant {
		qry += " AND software_development_relevant = ?"
		args = append(args, app.SoftwareDevelopmentRelevant)
	}
	if app.CloudOnly {
		qry += " AND cloud_only = ?"
		args = append(args, app.CloudOnly)
	}
	if app.PhysicalSecurityOnly {
		qry += " AND physical_security_only = ?"
		args = append(args, app.PhysicalSecurityOnly)
	}
	if app.PersonalSecurityOnly {
		qry += " AND personal_security_only = ?"
		args = append(args, app.PersonalSecurityOnly)
	}

	rows, err := r.db.QueryContext(ctx, qry, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ctrls []check.Control
	for rows.Next() {
		ctrl, err := scanForControl(rows)
		if err != nil {
			log.Printf("Scanning SQL row: %v\n", err)
			continue
		}
		ctrls = append(ctrls, ctrl)
	}

	return ctrls, nil
}

func (r *repository) SaveForApplication(ctx context.Context, app *application.Application, ctrls []check.Control) error {

	qry := `insert into APP_CONTROLS (APP_ID, CHECK_ID) values (?, ?)`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, qry)
	if err != nil {
		return err
	}

	for _, c := range ctrls {
		_, err := stmt.ExecContext(ctx, app.ID, c.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
