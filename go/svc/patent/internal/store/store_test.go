package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"sort"
	"testing"
	"time"

	"github.com/patential/go/pkg/testutil"
	"github.com/sirupsen/logrus/hooks/test"
)

var logger, _ = test.NewNullLogger()

const testDBSchemaPath = "./testdata/schema.sql"

func TestListPatents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := testutil.NewSqliteClient(ctx, testDBSchemaPath)
	if err != nil {
		t.Fatalf("Failed to prepare test DB: %v", err)
	}
	defer db.Close()

	// moment is a "frozen" time for the current test context.
	moment := time.Now()
	future := moment.Add(time.Duration(1))

	patents := []patentRow{
		patentRow{
			applicationNumber: "application-number",
			applicationKind:   "application-kind",
			grantDate:         moment,
		},
		patentRow{
			applicationNumber: "application-number-2",
			applicationKind:   "application-kind",
			grantDate:         future,
		},
	}

	for _, p := range patents {
		_, err := insertPatent(db, p)
		if err != nil {
			t.Fatalf("Failed to insert patent: %v", err)
		}
	}

	c := NewClient(logger, db)
	list, err := c.ListPatents(ctx, Pagination{})
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
		return
	}

	if len(list) != len(patents) {
		t.Errorf(
			"Wrong size of patents list returned (got=%d, expected=%d)",
			len(list), len(patents),
		)
		return
	}

	sorted := sort.SliceIsSorted(
		list,
		func(i, j int) bool { return list[j].GrantDate.Before(list[i].GrantDate) },
	)
	if !sorted {
		t.Errorf(
			"Wrong order of patents list returned (got=%v)",
			list,
		)
	}
}

func insertPatent(db *sql.DB, row patentRow) (string, error) {
	sqlStatement := `
		INSERT INTO patents (
			id,
			application_number,
			application_kind,
			grant_date
		)
		VALUES ($1, $2, $3, $4)
	`
	st, err := db.Prepare(sqlStatement)
	if err != nil {
		return "", fmt.Errorf("failed to prepare DB statement: %v", err)
	}

	id := uuid.New().String()

	_, err = st.Exec(
		id,
		row.applicationNumber,
		row.applicationKind,
		row.grantDate,
	)
	if err != nil {
		return "", fmt.Errorf("failed to exec DB statement: %v", err)
	}

	return id, nil
}
