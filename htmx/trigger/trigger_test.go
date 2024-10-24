package trigger_test

import (
	"testing"

	"github.com/patrickward/gtml/htmx/trigger"
)

func TestTriggers_Encode(t *testing.T) {
	triggers := trigger.NewTriggers()

	tests := []struct {
		name     string
		triggers map[string]*trigger.Trigger
		expected string
		wantErr  bool
	}{
		{
			name:     "empty",
			triggers: map[string]*trigger.Trigger{},
			expected: "",
			wantErr:  false,
		},
		{
			name: "with simple triggers",
			triggers: map[string]*trigger.Trigger{
				"test":  trigger.NewTrigger("test", nil),
				"test2": trigger.NewTrigger("test2", nil),
			},
			expected: "{\"test\":\"\",\"test2\":\"\"}",
			wantErr:  false,
		},
		{
			name: "with string value triggers",
			triggers: map[string]*trigger.Trigger{
				"test":  trigger.NewTrigger("test", "value"),
				"test2": trigger.NewTrigger("test2", "value"),
			},
			expected: "{\"test\":\"value\",\"test2\":\"value\"}",
			wantErr:  false,
		},
		{
			name: "with complex json triggers",
			triggers: map[string]*trigger.Trigger{
				"test":  trigger.NewTrigger("test", map[string]any{"key": "value"}),
				"test2": trigger.NewTrigger("test2", map[string]any{"key": "value"}),
				"test3": trigger.NewTrigger("test3", map[string]any{"key": 3}),
				"test4": trigger.NewTrigger("test4", map[string]any{"key": true}),
			},
			expected: "{\"test\":{\"key\":\"value\"},\"test2\":{\"key\":\"value\"},\"test3\":{\"key\":3},\"test4\":{\"key\":true}}",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := triggers.Encode(tt.triggers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Triggers.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("Triggers.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
