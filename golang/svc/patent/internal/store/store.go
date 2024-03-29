package store

import (
	"context"
	"fmt"
	"time"

	"database/sql"
	"github.com/google/uuid"
	lg "github.com/patential/golang/pkg/log"
)

// Patent represents publication of grant patent.
type Patent struct {
	ID                string
	ApplicationNumber string
	ApplicationKind   string
	GrantDate         time.Time
}

// Pagination defines pagination attributes.
type Pagination struct {
	Offset int32
	Limit  int32
}

// DB defines methods needed to interact with Postgres DB.
type DB interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Client represents store client.
type Client struct {
	logger lg.Logger
	db     DB
}

// patentRow represents patent DB row.
type patentRow struct {
	id                string
	applicationNumber string
	applicationKind   string
	grantDate         time.Time
}

// defaultLimit is the default pagination limit.
const defaultLimit = 25

// NewClient creates a new instance of Client.
func NewClient(l lg.Logger, db DB) Client {
	return Client{logger: l, db: db}
}

// ListPatents retrieves patents list with specified pagination parameters.
// The collection returned is ordered by `grant_date` timestamp.
func (c Client) ListPatents(ctx context.Context, pagination Pagination) ([]Patent, error) {
	limit := pagination.Limit

	if limit == 0 {
		limit = defaultLimit
	}

	rows, err := c.db.QueryContext(
		ctx,
		`
			SELECT
				id,
				application_number,
				application_kind,
				grant_date
			FROM patents
			ORDER BY grant_date DESC
			LIMIT $1
			OFFSET $2
		`,
		limit,
		pagination.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var patents []Patent

	for rows.Next() {
		var r patentRow

		err := rows.Scan(
			&r.id,
			&r.applicationNumber,
			&r.applicationKind,
			&r.grantDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error copying row values to destination: %v", err)
		}

		patents = append(
			patents,
			Patent{
				ID:                r.id,
				ApplicationNumber: r.applicationNumber,
				ApplicationKind:   r.applicationKind,
				GrantDate:         r.grantDate,
			},
		)
	}

	return patents, nil
}

// InsertPatent inserts single patent.
// It returns generated ID in case of success.
func (c Client) InsertPatent(ctx context.Context, patent Patent) (string, error) {
	id := uuid.New().String()

	row := patentRow{
		id:                id,
		applicationNumber: patent.ApplicationNumber,
		applicationKind:   patent.ApplicationKind,
		grantDate:         patent.GrantDate,
	}

	_, err := c.db.ExecContext(
		ctx,
		`
			INSERT INTO patents (
				id,
				application_number,
				application_kind,
				grant_date
			)
			VALUES ($1, $2, $3, $4)
		`,
		row.id,
		row.applicationNumber,
		row.applicationKind,
		row.grantDate,
	)
	if err != nil {
		return "", fmt.Errorf("failed to exec DB statement: %v", err)
	}

	return id, nil
}
