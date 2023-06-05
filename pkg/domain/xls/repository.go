package xls

import (
	"context"
	"fmt"
	"strings"

	"github.com/mehix/sec-checklist/pkg/domain"
	"github.com/xuri/excelize/v2"
)

type repository struct {
	fPath     string
	sheetName string
}

func NewRepository(fPath, sheetName string) domain.Reader {
	return &repository{fPath: fPath, sheetName: sheetName}
}

func (r *repository) FetchByID(_ context.Context, id string) (domain.Control, error) {
	return domain.Control{}, fmt.Errorf("not implemented")
}

func (r *repository) FetchByType(t string) ([]domain.Control, error) {
	all, err := r.FetchAll()
	if err != nil {
		return nil, err
	}

	var byType []domain.Control
	for _, c := range all {
		if c.Type == t {
			byType = append(byType, c)
		}
	}

	return byType, nil
}

func (r *repository) FetchAll() ([]domain.Control, error) {
	f, err := excelize.OpenFile(r.fPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if r.sheetName == "" {
		r.sheetName = f.GetSheetList()[0]
	}
	fmt.Printf("Reading sheet '%s'\n", r.sheetName)
	rows, err := f.Rows(r.sheetName)
	if err != nil {
		return nil, err
	}

	counter := 0
	controls := make([]domain.Control, 0)
	for rows.Next() {
		if counter == 0 {
			//header row
			counter++
			continue
		}

		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if ctrl := fromRowData(row); ctrl.ID != "" {
			controls = append(controls, ctrl)
		}
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}

	return controls, nil
}

func fromRowData(row []string) domain.Control {
	// some rows may have less columns filled with data
	cols := make([]string, 30)
	copy(cols, row)

	entry := domain.Control{
		Type: strings.TrimSpace(cols[0]),
		ID:   strings.TrimSpace(cols[1]),
		Name: strings.TrimSpace(cols[2]),
	}

	entry.Description = strings.TrimSpace(cols[3])
	entry.C = strings.TrimSpace(cols[4])
	entry.I = strings.TrimSpace(cols[5])
	entry.A = strings.TrimSpace(cols[6])
	entry.T = strings.TrimSpace(cols[7])
	entry.PD = strings.TrimSpace(cols[8])
	entry.NSI = strings.TrimSpace(cols[9])
	entry.SESE = strings.TrimSpace(cols[10])
	entry.OTCL = strings.TrimSpace(cols[11])
	entry.CSRDirection = strings.TrimSpace(cols[12])
	entry.SPSA = strings.TrimSpace(cols[13])
	entry.SPSAUnique = strings.TrimSpace(cols[14])
	entry.GDPR = strings.TrimSpace(cols[15]) == "GDPR"
	entry.GDPRUnique = strings.TrimSpace(cols[16]) == "GDPR unique"
	entry.ExternalSupplier = strings.TrimSpace(cols[17]) != ""
	entry.AssetType = strings.TrimSpace(cols[18])
	entry.OperationalCapability = strings.TrimSpace(cols[19])
	entry.PartOfGISR = strings.ToLower(strings.TrimSpace(cols[20])) == "yes"
	entry.LastUpdated = strings.TrimSpace(cols[21])
	entry.OldID = strings.TrimSpace(cols[22])

	entry.OnlyHandledCentrally = strings.TrimSpace(cols[23]) == "yes"
	entry.HandledCentrallyBy = strings.TrimSpace(cols[24])
	entry.ExcludeForExternalSupplier = strings.TrimSpace(cols[25]) == "yes"
	entry.SoftwareDevelopmentRelevant = strings.TrimSpace(cols[26]) == "yes"
	entry.CloudOnly = strings.TrimSpace(cols[27]) == "yes"
	entry.PhysicalSecurityOnly = strings.TrimSpace(cols[28]) == "yes"
	entry.PersonalSecurityOnly = strings.TrimSpace(cols[29]) == "yes"

	return entry
}
