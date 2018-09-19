package user

import (
	"fmt"

	pb "github.com/GoingFast/gotrains/user/protobuf"
	"google.golang.org/grpc"
)

func Client() (*grpc.ClientConn, pb.UserServiceClient) {
	conn, err := grpc.Dial(":50052", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	c := pb.NewUserServiceClient(conn)
	return conn, c
}
