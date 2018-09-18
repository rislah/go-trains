package train

import (
	"context"
	"database/sql"
	"fmt"

	pb "github.com/GoingFast/gotrains/train/protobuf"
	"github.com/GoingFast/gotrains/util/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	service struct {
		store  trainStore
		log    logger.Log
		logrus *logrus.Logger
	}
)

func NewTrainService(db *sql.DB, l logger.Log, lgrus *logrus.Logger) service {
	s := newTrainStore(db)
	return service{
		store:  s,
		log:    l,
		logrus: lgrus,
	}
}

func (s service) CreateTrain(ctx context.Context, req *pb.Train) (*pb.CreateTrainResponse, error) {

	ex, err := s.store.trainExists(req.Brandname)
	if err != nil {
		s.logrus.Debug(err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	if ex {
		return nil, status.Errorf(codes.AlreadyExists, "train already exists")
	}

	err = s.store.createTrain(req)
	if err != nil {
		s.logrus.Debug(err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	go s.log.AuditLog(ctx, fmt.Sprintf("created new train with the name %s", req.Brandname))

	return &pb.CreateTrainResponse{Msg: "success"}, nil
}
