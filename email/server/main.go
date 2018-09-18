package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/GoingFast/gotrains/email"
	pb "github.com/GoingFast/gotrains/email/protobuf"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

const (
	GRPCPort = ":50053"
)

func main() {
	ln, err := net.Listen("tcp", GRPCPort)
	if err != nil {
		log.Fatalf("grpc tcp: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterEmailServiceServer(server, email.NewEmailService())
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var g group.Group
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
		g.Add(func() error {
			fmt.Println("listening...")
			return server.Serve(ln)
		}, func(error) {
			server.GracefulStop()
		})
	}

	log.Fatal("exit ", g.Run())
}
