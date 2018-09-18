package user

import (
	"context"
	"database/sql"

	"github.com/GoingFast/gotrains/email"
	pbemail "github.com/GoingFast/gotrains/email/protobuf"
	pb "github.com/GoingFast/gotrains/user/protobuf"
	"github.com/GoingFast/gotrains/util/auth"
	"github.com/GoingFast/gotrains/util/logger"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	store  userStore
	log    logger.Log
	logrus *logrus.Logger
}

// NewUserService returns a new instance on service
func NewUserService(db *sql.DB, l logger.Log, lgrus *logrus.Logger) service {
	s := newUserStore(db)
	return service{
		store:  s,
		log:    l,
		logrus: lgrus,
	}
}

func (s service) CreateUser(ctx context.Context, req *pb.User) (*pb.CreateUserResponse, error) {
	_, ex, err := s.store.findUserByEmail(req.Email)
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "user")
	}

	if ex {
		return nil, status.Errorf(codes.AlreadyExists, "user with the requested email already exists")
	}

	h, err := bcrypt.GenerateFromPassword([]byte(req.Password), 3)
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "user")
	}

	vid, err := ksuid.NewRandom()
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "user")
	}

	req.Password = string(h)

	err = s.store.createUser(req, vid.String())
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "user")
	}

	conn, client := email.Client()
	email := &pbemail.Email{
		To:             req.Email,
		Verificationid: vid.String(),
	}

	go func() {
		defer conn.Close()
		_, err := client.SendEmailVerification(context.Background(), email)
		if err != nil {
			s.logrus.Debug(err)
		}
	}()

	return &pb.CreateUserResponse{Msg: "success"}, nil
}

func (s service) Login(ctx context.Context, req *pb.User) (*pb.LoginResponse, error) {
	u, ex, err := s.store.findUserByUsername(req.Username)
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "user")
	}

	if !ex {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	if ex && !u.Verified {
		return nil, status.Errorf(codes.Unauthenticated, "account is waiting to be verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	token, err := auth.EncodeJWT(u)
	if err != nil {
		s.logrus.Debug(err)
		return nil, s.log.InternalError(err, "user")
	}

	return &pb.LoginResponse{Token: token}, nil
}
