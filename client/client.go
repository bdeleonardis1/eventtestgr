package client

import (
	"log"

	"google.golang.org/grpc"

	pb "github.com/bdeleonardis1/eventtestgr/events"
)

func GetConnection() pb.EventsClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial("localhost:11111", opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	// TODO: figure out what to do about client defer conn.Close()

	client := pb.NewEventsClient(conn)

	return client
}