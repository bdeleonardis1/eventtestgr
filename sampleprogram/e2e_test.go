package main_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/bdeleonardis1/eventtestgr/events"
	"github.com/bdeleonardis1/eventtestgr/eventtest"
)

const (
	expectedBase = "Enter a number 1-5: "
)

func TestParity(t *testing.T) {
	server := eventtest.StartListening()
	defer server.GracefulStop()

	testCases := []struct {
		input          string
		expected       string
		expectedEvents []*events.Event
	}{
		{
			input:    "1",
			expected: "1 is an odd number",
			expectedEvents: []*events.Event{
				eventtest.NewEvent("1Optimization"),
			},
		},
		{
			input:    "2",
			expected: "2 is an even number",
			expectedEvents: []*events.Event{
				eventtest.NewEvent("OptimizedSingleDigit"),
			},
		},
		{
			input:    "11",
			expected: "11 is an odd number",
			expectedEvents: []*events.Event{
				eventtest.NewEvent("convertToNumber"), eventtest.NewEvent("Modding"), eventtest.NewEvent("TheVeryEnd"),
			},
		},
		{
			input:    "-3",
			expected: "-3 is an odd number",
			expectedEvents: []*events.Event{
				eventtest.NewEvent("convertToNumber"), eventtest.NewEvent("OptimizedNegativeSingleDigit"),
			},
		},
		{
			input:    "-4",
			expected: "-4 is an even number",
			expectedEvents: []*events.Event{
				eventtest.NewEvent("convertToNumber"), eventtest.NewEvent("OptimizedNegativeSingleDigit"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			eventtest.ClearEvents()

			cmd := exec.Command("./sampleprogram")
			cmd.Stdin = strings.NewReader(tc.input)
			out, err := cmd.Output()
			if err != nil {
				t.Fatal(err)
			}
			outString := strings.TrimSpace(string(out))

			if outString != expectedBase+tc.expected {
				t.Errorf("expected '%v', but got '%v'", expectedBase+tc.expected, outString)
			}

			eventtest.ExpectExactEvents(t, tc.expectedEvents)
		})
	}
}

func TestExpectEventsDemo(t *testing.T) {
	server := eventtest.StartListening()
	defer server.GracefulStop()

	eventtest.ClearEvents()

	cmd := exec.Command("./sampleprogram")
	cmd.Stdin = strings.NewReader("19")
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	eventtest.ExpectEvents(t, []*events.Event{eventtest.NewEvent("convertToNumber"), eventtest.NewEvent("TheVeryEnd")}, eventtest.Ordered)
	eventtest.ExpectEvents(t, []*events.Event{eventtest.NewEvent("TheVeryEnd"), eventtest.NewEvent("convertToNumber")}, eventtest.Unordered)

	eventtest.UnexpectedEvents(t, []*events.Event{eventtest.NewEvent("1Optimization"), eventtest.NewEvent("OptimizedNegativeSingleDigit")})

	// should fail.
	eventtest.ExpectEvents(t, []*events.Event{eventtest.NewEvent("TheVeryEnd"), eventtest.NewEvent("convertToNumber")}, eventtest.Ordered)
}
