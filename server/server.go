package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/bdeleonardis1/eventtestgr/events"
)

type eventsServer struct {
	events []*pb.Event
	pb.UnimplementedEventsServer
}

func newEventsServer() *eventsServer {
	return &eventsServer{
		events: make([]*pb.Event, 0),
	}
}

func getEmpty() *pb.Empty {
	return &pb.Empty{}
}

func (es *eventsServer) EmitEvent(ctx context.Context, event *pb.Event) (*pb.Empty, error) {
	es.events = append(es.events, event)

	return getEmpty(), nil
}

func (es *eventsServer) GetEvents(ctx context.Context, _ *pb.Empty) (*pb.EventList, error) {
	return &pb.EventList{ Events: es.events}, nil
}

func (es *eventsServer) ClearEvents(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	es.events = make([]*pb.Event, 0)

	return getEmpty(), nil
}

func StartServer() {
	lis, err := net.Listen("tcp", "localhost:11111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEventsServer(grpcServer, newEventsServer())
	go grpcServer.Serve(lis)
}
