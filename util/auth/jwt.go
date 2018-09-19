package auth

import (
	"context"
	"strings"
	"time"

	pb "github.com/GoingFast/gotrains/user/protobuf"
	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// jwt token expiration time
	expirationTime = time.Now().Add(6 * time.Hour).Unix()
	// secret which is used to decode and encode the token
	secret = "changeme"
)

const JWTClaims contextKey = "JWTClaims"

type (
	Claims struct {
		Username string `json:"username,omitempty"`
		Role     string `json:"role,omitempty"`
		jwt.StandardClaims
	}

	contextKey string
)

func EncodeJWT(u pb.User) (string, error) {
	claims := Claims{
		u.Username,
		u.Role,
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func decodeJWT(t string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(t, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "bad authorization token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok && !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "bad authorization token")
	}
	return claims, nil
}

func extractJWT(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get("authorization")
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "bad authorization token schema")
	}

	if strings.ToLower(splits[0]) != "bearer" {
		return "", status.Errorf(codes.Unauthenticated, "bad authorization token schema")
	}

	return splits[1], nil
}

func GetJWTClaims(ctx context.Context) (*Claims, error) {
	token, err := extractJWT(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := decodeJWT(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func Middleware(methods ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error

		for _, method := range methods {
			if method == info.FullMethod {
				newCtx, err = func(ctx context.Context) (context.Context, error) {
					claims, err := GetJWTClaims(ctx)
					if err != nil {
						return nil, err
					}
					newCtx = context.WithValue(ctx, JWTClaims, claims)
					return newCtx, nil
				}(ctx)
			}
		}

		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

func CheckRole(ctx context.Context, roles ...string) error {
	claims, err := GetJWTClaims(ctx)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if role == claims.Role {
			return status.Errorf(codes.Unauthenticated, "you are not privileged to do that")
		}
	}

	return nil
}
