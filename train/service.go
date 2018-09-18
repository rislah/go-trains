package train

import (
	"context"
	"database/sql"
	"fmt"

	pb "github.com/GoingFast/gotrains/train/protobuf"
	"github.com/GoingFast/gotrains/util/auth"
	"github.com/GoingFast/gotrains/util/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	store  trainStore
	log    logger.Log
	logrus *logrus.Logger
}

// NewTrainService returns a new instance of service
func NewTrainService(db *sql.DB, l logger.Log, lgrus *logrus.Logger) service {
	s := newTrainStore(db)
	return service{
		store:  s,
		log:    l,
		logrus: lgrus,
	}
}

func (s service) CreateTrain(ctx context.Context, req *pb.Train) (*pb.CreateTrainResponse, error) {
	err := auth.CheckRole(ctx, []string{"customer"}...)
	if err != nil {
		return nil, err
	}

	ex, err := s.store.trainExists(req.Brandname)
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "train")
	}

	if ex {
		return nil, status.Errorf(codes.AlreadyExists, "train already exists")
	}

	err = s.store.createTrain(req)
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "train")
	}

	go s.log.AuditLog(ctx, fmt.Sprintf("created new train: %s", req.Brandname))

	return &pb.CreateTrainResponse{Msg: "success"}, nil
}

func (s service) GetTrains(ctx context.Context, req *pb.Empty) (*pb.GetTrainsResponse, error) {
	err := auth.CheckRole(ctx, []string{"customer"}...)
	if err != nil {
		return nil, err
	}

	t, err := s.store.getTrains()
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "train")
	}

	if len(t) <= 0 {
		return nil, status.Errorf(codes.NotFound, "couldn't find any trains")
	}

	return &pb.GetTrainsResponse{Trains: t}, nil
}

func (s service) CreateRoute(ctx context.Context, req *pb.Route) (*pb.CreateRouteResponse, error) {
	err := auth.CheckRole(ctx, []string{"customer"}...)
	if err != nil {
		return nil, err
	}

	ex, err := s.store.trainExists(req.Brandname)
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "train")
	}

	if !ex {
		return nil, status.Errorf(codes.NotFound, "train with the requested requested name not found")
	}

	err = s.store.createRoute(req)
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "train")
	}

	go s.log.AuditLog(ctx, fmt.Sprintf("created a new path: from %v, to %v, train %v", req.From, req.To, req.Brandname))

	return &pb.CreateRouteResponse{Msg: "success"}, nil
}
