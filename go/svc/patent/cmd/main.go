package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net"
	"time"

	_ "github.com/lib/pq"
	lg "github.com/patential/go/pkg/log"
	"github.com/patential/go/svc/patent/internal/service"
	"github.com/patential/go/svc/patent/internal/store"
	pb "github.com/patential/go/svc/patent/proto/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// sqlDriverName is driver used to connect to Patents DB.
	sqlDriverName = "postgres"
)

var (
	dbHost     = flag.String("db-host", "", "Hostname for Postgresql instance.")
	dbPort     = flag.Int("db-port", 0, "Port number for Postgresql instance.")
	dbUser     = flag.String("db-user", "", "Username for Postgresql instance.")
	dbPassword = flag.String("db-password", "", "Password for Postgresql instance.")
	dbDatabase = flag.String("db-name", "", "Patents database name on Postgresql instance.")

	port = flag.Int64("port", 0, "Server port")
)

var logger = lg.NewLogger()

func main() {
	if err := run(); err != nil {
		logger.Fatalf("Error running Patent Service: %v", err)
	}
}

func run() error {
	flag.Parse()

	if err := areAllFlagsSet(); err != nil {
		return err
	}

	logger.Infoln("Starting Patent Service...")

	dbURI := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		*dbHost,
		*dbPort,
		*dbUser,
		*dbPassword,
		*dbDatabase,
	)

	logger.Infoln("Connecting to Patents DB...")
	db, err := sql.Open(sqlDriverName, dbURI)
	if err != nil {
		return fmt.Errorf("failed to connect to Patents DB: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)

	sc := store.NewClient(logger, db)
	ps := service.New(logger, sc)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterServiceServer(grpcServer, ps)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	logger.Infof("Service listening on port: %d", *port)
	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

func areAllFlagsSet() error {
	passed := map[string]*flag.Flag{}
	flag.Visit(func(f *flag.Flag) {
		passed[f.Name] = f
	})

	all := []*flag.Flag{}
	flag.VisitAll(func(f *flag.Flag) {
		all = append(all, f)
	})
	if len(passed) != len(all) {
		msg := ""
		for _, f := range all {
			msg = msg + f.Name + ", "
		}
		return fmt.Errorf("expected %d flags but got only %d, required flags are: %s", len(all), len(passed), msg)
	}
	return nil
}
