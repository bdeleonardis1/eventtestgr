package main_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/bdeleonardis1/eventtestgr/api"
	"github.com/bdeleonardis1/eventtestgr/events"
)

const (
	expectedBase = "Enter a number 1-5: "
)

func TestParity(t *testing.T) {
	server := api.StartListening()
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
				api.NewEvent("1Optimization"),
			},
		},
		{
			input:    "2",
			expected: "2 is an even number",
			expectedEvents: []*events.Event{
				api.NewEvent("OptimizedSingleDigit"),
			},
		},
		{
			input:    "11",
			expected: "11 is an odd number",
			expectedEvents: []*events.Event{
				api.NewEvent("convertToNumber"), api.NewEvent("Modding"), api.NewEvent("TheVeryEnd"),
			},
		},
		{
			input:    "-3",
			expected: "-3 is an odd number",
			expectedEvents: []*events.Event{
				api.NewEvent("convertToNumber"), api.NewEvent("OptimizedNegativeSingleDigit"),
			},
		},
		{
			input:    "-4",
			expected: "-4 is an even number",
			expectedEvents: []*events.Event{
				api.NewEvent("convertToNumber"), api.NewEvent("OptimizedNegativeSingleDigit"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			api.ClearEvents()

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

			api.ExpectExactEvents(t, tc.expectedEvents)
		})
	}
}

func TestExpectEventsDemo(t *testing.T) {
	server := api.StartListening()
	defer server.GracefulStop()

	api.ClearEvents()

	cmd := exec.Command("./sampleprogram")
	cmd.Stdin = strings.NewReader("19")
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	api.ExpectEvents(t, []*events.Event{api.NewEvent("convertToNumber"), api.NewEvent("TheVeryEnd")}, api.Ordered)
	api.ExpectEvents(t, []*events.Event{api.NewEvent("TheVeryEnd"), api.NewEvent("convertToNumber")}, api.Unordered)

	api.UnexpectedEvents(t, []*events.Event{api.NewEvent("1Optimization"), api.NewEvent("OptimizedNegativeSingleDigit")})

	// should fail.
	api.ExpectEvents(t, []*events.Event{api.NewEvent("TheVeryEnd"), api.NewEvent("convertToNumber")}, api.Ordered)
}
