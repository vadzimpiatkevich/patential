package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/go-cmp/cmp"
	"github.com/patential/go/pkg/testutil"
	"github.com/patential/go/svc/patent/internal/store"
	pb "github.com/patential/go/svc/patent/proto/service"
	lgtest "github.com/sirupsen/logrus/hooks/test"
)

var logger, _ = lgtest.NewNullLogger()

const testDBSchemaPath = "./testdata/schema.sql"

func TestService_ListPatents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := testutil.NewSqliteClient(ctx, testDBSchemaPath)
	if err != nil {
		t.Fatalf("Failed to prepare test DB: %v", err)
	}
	defer db.Close()

	storeClient := store.NewClient(logger, db)

	// moment is a "frozen" time for the current test context.
	moment := time.Now()

	// patents is the list of patent instances to insert.
	patents := []store.Patent{
		store.Patent{
			ApplicationNumber: "application-number",
			ApplicationKind:   "application-kind",
			GrantDate:         moment.Add(time.Duration(1)),
		},
		store.Patent{
			ApplicationNumber: "application-number-2",
			ApplicationKind:   "application-kind-2",
			GrantDate:         moment,
		},
	}

	for _, p := range patents {
		_, err := storeClient.InsertPatent(ctx, p)
		if err != nil {
			t.Fatalf("Failed to insert patent: %v", err)
		}
	}

	service := New(logger, storeClient)

	t.Run("with pagination params", func(t *testing.T) {
		const (
			limit  = 1
			offset = 1
		)

		resp, err := service.ListPatents(
			ctx,
			&pb.ListPatentsRequest{
				Pagination: &pb.Pagination{Limit: limit, Offset: offset},
			},
		)
		if err != nil {
			t.Errorf("Unexpected error returned: %v", err)
			return
		}
		if len(resp.Patents) != limit {
			t.Errorf(
				"Wrong size of patents list returned (got=%d, expected=%d)",
				len(resp.Patents), limit,
			)
			return
		}

		for i, p := range resp.Patents {
			gdt, err := ptypes.Timestamp(p.GrantDate)
			if err != nil {
				t.Fatalf(
					"Failed to convert grant date (grantDate=%v)", p.GrantDate,
				)
			}

			// Map patent instance from Proto to Store.
			got := store.Patent{
				ApplicationNumber: p.ApplicationNumber,
				ApplicationKind:   p.ApplicationKind,
				GrantDate:         gdt,
			}
			want := patents[i+offset]

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf(
					"Unexpected patent instance returned at index %d (diff=%s)", i, diff,
				)
			}
		}
	})
}
