package eventtest

import (
	"context"

	"github.com/bdeleonardis1/eventtestgr/client"
	pb "github.com/bdeleonardis1/eventtestgr/events"
	"github.com/bdeleonardis1/eventtestgr/server"
)

type StoppableServer interface {
	GracefulStop()
}

func StartListening() StoppableServer {
	return server.StartServer()
}

func EmitEvent(name string) error {
	c := client.GetConnection()
	_, err := c.EmitEvent(context.Background(), &pb.Event{Name: name})
	return err
}

func ClearEvents() error {
	c := client.GetConnection()
	_, err := c.ClearEvents(context.Background(), getEmpty())
	return err
}

func GetEvents() ([]*pb.Event, error) {
	c := client.GetConnection()
	eventList, err := c.GetEvents(context.Background(), getEmpty())
	if err != nil {
		return nil, err
	}
	return eventList.Events, nil
}

func GetNameList(events []*pb.Event) []string {
	nameList := make([]string, len(events))
	for i, e := range events {
		nameList[i] = e.Name
	}
	return nameList
}

func getEmpty() *pb.Empty {
	return &pb.Empty{}
}


func NewEvent(name string) *pb.Event {
	return &pb.Event{
		Name: name,
	}
}