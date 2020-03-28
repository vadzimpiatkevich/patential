package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// sqliteDriver is driver used to connect to test DB.
	sqliteDriver = "sqlite3"
	// sqliteTempfileName is filename given to tempfile used by test DB.
	sqliteTempfileName = "sqlite3.db"
)

// NewSqliteClient creates test database and returns Sqlite client instance.
func NewSqliteClient(ctx context.Context, schemaPath string) (*sql.DB, error) {
	// Create tempfile used by test DB.
	tmpfile, err := ioutil.TempFile("", sqliteTempfileName)
	if err != nil {
		return nil, fmt.Errorf("could not create temporal file used by DB: %v", err)
	}

	// Open connection to DB.
	db, err := sql.Open(sqliteDriver, tmpfile.Name())
	if err != nil {
		return nil, fmt.Errorf("error opening DB connection: %v", err)
	}

	// Read "raw" bytes from schema file.
	bytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error reading DB statements: %v", err)
	}

	// schema is plain DB schema text.
	schema := string(bytes)
	// queries is schema separated by delimiter and represents an array of DB
	// queries to execute.
	queries := strings.Split(schema, ";\n")

	// Iterate over schema queries and execute.
	for _, q := range queries[:len(queries)-1] {
		st, err := db.Prepare(q)
		if err != nil {
			defer db.Close()
			return nil, fmt.Errorf("error preparing DB statement from query: %v", err)
		}
		defer st.Close()

		if _, err := st.Exec(); err != nil {
			defer db.Close()
			return nil, fmt.Errorf("error executing query: %v", err)
		}
	}

	return db, nil
}
