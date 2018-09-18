package main

import (
	"context"
	"database/sql"
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"os"

	"github.com/GoingFast/gotrains/train"
	pb "github.com/GoingFast/gotrains/train/protobuf"
	"github.com/GoingFast/gotrains/util/auth"
	"github.com/GoingFast/gotrains/util/logger"
	sentry "github.com/getsentry/raven-go"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"github.com/oklog/oklog/pkg/group"
	"github.com/olivere/elastic"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	HTTPPort = ":8080"
	GRPCPort = ":50051"
)

func main() {
	var log *logrus.Logger
	{
		log = logrus.New()
	}
	log.AddHook(filename.NewHook())

	if os.Getenv("debug") == "1" || os.Getenv("debug") == "" {
		log.SetLevel(logrus.DebugLevel)
	}

	ln, err := net.Listen("tcp", GRPCPort)
	if err != nil {
		stdlog.Fatalf("grpc tcp: %v", err)
	}

	uri := fmt.Sprintf("postgres://yeye:mysecretpassword@localhost:5432/testing?sslmode=disable")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		stdlog.Fatalf("postgres: %v", err)
	}

	_, err = db.Exec("SELECT 1")
	if err != nil {
		stdlog.Fatalf("postgres test query: %v", err)
	}

	defer db.Close()

	s, err := sentry.NewWithTags("http://b6c78f9743c944f0bbcf26c4c2c0b797:6a3584ace43b428191ffc7a16e81ff6e@localhost:9000/1", map[string]string{"service": "train"})
	if err != nil {
		stdlog.Fatalf("sentry: %v", err)
	}

	defer s.Close()

	e, err := elastic.NewClient()
	if err != nil {
		stdlog.Fatalf("elasticsearch: %v", err)
	}

	// endpoints that should require authentication
	// requires full protobuf method name
	server := grpc.NewServer(grpc.UnaryInterceptor(auth.Middleware([]string{
		"/train.TrainService/CreateTrain",
		"/train.TrainService/CreateRoute",
		"/train.TrainService/GetTrains",
	}...)))

	var logsvc logger.Log
	{
		logsvc = logger.NewService(s, db, e)
	}

	err = logsvc.SetupElasticIndexes()
	if err != nil {
		stdlog.Fatalf("elastic indexes; %v", err)
	}

	trainSvc := train.NewTrainService(db, logsvc, log)

	pb.RegisterTrainServiceServer(server, trainSvc)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	conn, err := grpc.Dial(GRPCPort, grpc.WithInsecure())
	if err != nil {
		stdlog.Fatalf("grpc dial: %v", err)
	}

	err = pb.RegisterTrainServiceHandler(ctx, mux, conn)
	if err != nil {
		stdlog.Fatalf("grpc register handler: %v", err)
	}

	var g group.Group
	{
		g.Add(func() error {
			return server.Serve(ln)
		}, func(error) {
			server.GracefulStop()
		})
	}
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			sigchan := make(chan os.Signal, 1)
			select {
			case sig := <-sigchan:
				return fmt.Errorf("received signal: %v", sig)
			case <-cancel:
				return nil
			}

		}, func(error) {
			close(cancel)
		})
	}
	{
		ln, err := net.Listen("tcp", HTTPPort)
		if err != nil {
			stdlog.Fatalf("http tcp: %v", err)
		}

		g.Add(func() error {
			fmt.Println("listening...")
			return http.Serve(ln, mux)
		}, func(error) {
			ln.Close()
			log.Fatalf("got an error: %v", err)
		})
	}
	stdlog.Fatal("exit ", g.Run())
}
