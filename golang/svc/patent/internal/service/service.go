package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"
	lg "github.com/patential/golang/pkg/log"
	"github.com/patential/golang/svc/patent/internal/store"
	pb "github.com/patential/golang/svc/patent/proto/service"
)

// Store is an interface needed to interact with a database.
type Store interface {
	ListPatents(ctx context.Context, pagination store.Pagination) ([]store.Patent, error)
}

// Service represents service.
type Service struct {
	logger lg.Logger
	store  Store
}

// New creates a new service.
func New(l lg.Logger, s Store) Service {
	return Service{
		logger: l,
		store:  s,
	}
}

// ListPatents returns list of patents with specified pagination parameters.
func (s Service) ListPatents(ctx context.Context, req *pb.ListPatentsRequest) (*pb.ListPatentsResponse, error) {
	s.logger.Infoln("ListPatents received")

	var pagination store.Pagination
	// Map pagination instance from Proto to Store.
	if rpg := req.GetPagination(); rpg != nil {
		pagination = store.Pagination{
			Limit:  rpg.Limit,
			Offset: rpg.Offset,
		}
	}

	patents, err := s.store.ListPatents(ctx, pagination)
	if err != nil {
		s.logger.Errorf("Failed to list patents: %v", err)
		return nil, status.Errorf(codes.Internal, "could not list patents")
	}

	response := new(pb.ListPatentsResponse)

	for _, p := range patents {
		// Convert time.Time to Proto timestamp.
		gdp, err := ptypes.TimestampProto(p.GrantDate)
		if err != nil {
			s.logger.Errorf("Failed to convert grant date: %v", err)
			return nil, status.Errorf(codes.Internal, "could not list patents")
		}

		// Map patent instance from Store to Proto.
		pbp := pb.Patent{
			Id:                p.ID,
			ApplicationNumber: p.ApplicationNumber,
			ApplicationKind:   p.ApplicationKind,
			GrantDate:         gdp,
		}

		response.Patents = append(response.Patents, &pbp)
	}

	return response, nil
}
