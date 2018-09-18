package email

import (
	"fmt"

	pb "github.com/GoingFast/gotrains/email/protobuf"
	"google.golang.org/grpc"
)

func Client() (*grpc.ClientConn, pb.EmailServiceClient) {
	conn, err := grpc.Dial(":50053", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	c := pb.NewEmailServiceClient(conn)
	return conn, c
}
