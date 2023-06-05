package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mehix/sec-checklist/pkg/domain"
)

type repository struct {
	db *sql.DB
}

func NewRepository(dsn string) domain.ReaderWriter {
	db, err := DbConnWithRetry(dbConn, 5, time.Second, time.Minute)(context.Background(), dsn)
	if err != nil {
		log.Fatalln("DB connection failed", err)
	}

	fmt.Println("Connected to database")
	return &repository{db: db}
}

func (r *repository) FetchAll() ([]domain.Control, error) {
	qry := "select * from V_CHECKS"

	rows, err := r.db.QueryContext(context.TODO(), qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ctrls []domain.Control
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
func (r *repository) FetchByType(t string) ([]domain.Control, error) { return nil, nil }

func (r *repository) FetchByID(ctx context.Context, id string) (domain.Control, error) {
	qry := "select * from V_CHECKS where ID = ?"

	row := r.db.QueryRowContext(ctx, qry, id)
	return scanForControl(row)
}

func (r *repository) SaveAll(ctx context.Context, all []domain.Control) (err error) {
	qry := "insert into CHECKS (ID, type, name, description,asset_type,last_update,old_id) values (?, ?, ?, ?, ?, ?, ?)"
	qryCiat := "insert CHECKS_CIAT (CHECK_ID, c, i, a, t) values (?, ?, ?, ?, ?)"
	qryFilters := `insert FILTERS (CHECK_ID, only_handle_centrally, 
		handled_centrally_by, excluded_for_external_supplier, 
		software_development_relevant, cloud_only, 
		physical_security_only, personal_security_only) 
		values (?, ?, ?, ?, ?, ?, ?, ?)`

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
	}

	return nil
}
