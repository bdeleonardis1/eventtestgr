package api

import (
	"context"
	"testing"

	"github.com/bdeleonardis1/eventtestgr/client"
	pb "github.com/bdeleonardis1/eventtestgr/events"
	"github.com/bdeleonardis1/eventtestgr/server"
)

type IsOrdered int

const (
	Ordered IsOrdered = iota
	Unordered
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

// ExpectExactEvents will error if the events that have been emitted
// do not exactly match the expectedEvents.
func ExpectExactEvents(t *testing.T, expectedEvents []*pb.Event) {
	t.Helper()

	actualEvents, err := GetEvents()
	if err != nil {
		t.Fatalf("error getting events: %v", err)
	}

	if len(actualEvents) != len(expectedEvents) {
		t.Fatalf("actual events: %v, not the same length as expected events: %v", GetNameList(actualEvents), GetNameList(expectedEvents))
	}

	for i, actualEvent := range actualEvents {
		if !eventEquals(actualEvent, expectedEvents[i]) {
			t.Errorf("the %vth actual event: %v, does not equal the expected event: %v", i, actualEvent, expectedEvents[i])
		}
	}
}

// ExpectEvents ensures that all expectedEvents have occurred. When ordered is
// Ordered, the expected events must occur in order in relation to each other.
// When ordered is Unordered, the events can occur in any order. This function
// ignores any events that are not in the expectedEvents list.
func ExpectEvents(t *testing.T, expectedEvents []*pb.Event, ordered IsOrdered) {
	t.Helper()

	actualEvents, err := GetEvents()

	if err != nil {
		t.Fatalf("error getting events: %v", err)
	}

	if ordered == Ordered {
		expectEventsOrdered(t, expectedEvents, actualEvents)
	} else {
		expectEventsUnordered(t, expectedEvents, actualEvents)
	}
}

func expectEventsOrdered(t *testing.T, expectedEvents, actualEvents []*pb.Event) {
	t.Helper()

	actualIdx := 0
	for _, expectedEvent := range expectedEvents {
		found := false
		for actualIdx < len(actualEvents) {
			if eventEquals(expectedEvent, actualEvents[actualIdx]) {
				found = true
				break
			}
			actualIdx += 1
		}
		if !found {
			t.Fatalf("could not find expected event: %v", expectedEvent)
		}
	}
}

func expectEventsUnordered(t *testing.T, expectedEvents, actualEvents []*pb.Event) {
	t.Helper()

	for _, expectedEvent := range expectedEvents {
		found := false
		for _, actualEvent := range actualEvents {
			if eventEquals(expectedEvent, actualEvent) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("could not find expected event: %v", expectedEvent)
		}
	}
}

// UnexpectedEvents if any of the events so far emitted
// are in the provided unexpectedEvents list.
func UnexpectedEvents(t *testing.T, unexpectedEvents []*pb.Event) {
	t.Helper()

	actualEvents, err := GetEvents()
	if err != nil {
		t.Fatal(err)
	}

	for _, unexpected := range unexpectedEvents {
		for _, actualEvent := range actualEvents {
			if eventEquals(unexpected, actualEvent) {
				t.Errorf("event: %v occurred even though it should not have", unexpected)
			}
		}
	}
}

func eventEquals(a, b *pb.Event) bool {
	return a.Name == b.Name
}

func NewEvent(name string) *pb.Event {
	return &pb.Event{
		Name: name,
	}
}