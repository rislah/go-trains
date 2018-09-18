package logger

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/GoingFast/gotrains/util/auth"
	sentry "github.com/getsentry/raven-go"
	"github.com/olivere/elastic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	auditLog struct {
		Name      string    `json:"name"`
		Action    string    `json:"action"`
		Timestamp time.Time `json:"timestamp"`
	}

	service struct {
		sentry *sentry.Client
		db     *sql.DB
		ec     *elastic.Client
	}

	// Log is a wrapper for logging methods
	Log interface {
		SetupElasticIndexes() error
		AuditLog(context.Context, string)
		InternalError(error, string) error
		// TODO
		// general logging
	}
)

// NewService returns a new instance of service
func NewService(s *sentry.Client, d *sql.DB, e *elastic.Client) *service {
	return &service{s, d, e}
}

// SetupElasticIndexes creates new indexes for our logging
func (s service) SetupElasticIndexes() error {
	ctx := context.Background()
	{
		exists, err := s.ec.IndexExists("auditlogs").Do(ctx)
		if err != nil {
			return err
		}

		if exists {
			return err
		}

		_, err = s.ec.CreateIndex("auditlogs").Do(ctx)
		if err != nil {
			return err
		}
	}
	{
		exists, err := s.ec.IndexExists("errorlogs").Do(ctx)
		if err != nil {
			return err
		}

		if exists {
			return err
		}

		_, err = s.ec.CreateIndex("generallogs").Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// AuditLog logs authorized users actions to Elasticsearch
// user info is retrieved from the context, user action needs to be passed in as the 2nd parameter
func (s service) AuditLog(ctx context.Context, a string) {
	claims, err := auth.GetJWTClaims(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	al := auditLog{
		Name:      claims.Username,
		Action:    a,
		Timestamp: time.Now().UTC(),
	}

	_, err = s.ec.Index().
		Index("auditlogs").
		Type("_doc").
		BodyJson(al).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
}

// InternalError takes unexpected errors, forwards them to Sentry and returns appropriate HTTP response
func (s service) InternalError(err error, svc string) error {
	_ = s.sentry.CaptureError(err, map[string]string{"service": svc})
	// TODO
	// html error page
	return status.Errorf(codes.Internal, "an internal server error has occured. Please contact technical support")
}
