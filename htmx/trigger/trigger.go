package trigger

import (
	"encoding/json"
)

// Trigger represents an HTMX trigger
// See the https://htmx.org/headers/hx-trigger/ documentation for more information.
type Trigger struct {
	name  string
	value any
}

// NewTrigger creates a new Trigger
func NewTrigger(name string, value any) *Trigger {
	return &Trigger{
		name:  name,
		value: value,
	}
}

// Triggers represents a collection of HTMX triggers
type Triggers struct {
	triggers    map[string]*Trigger
	afterSettle map[string]*Trigger
	afterSwap   map[string]*Trigger
}

// NewTriggers creates a new Triggers instance
func NewTriggers() *Triggers {
	return &Triggers{
		triggers:    make(map[string]*Trigger),
		afterSettle: make(map[string]*Trigger),
		afterSwap:   make(map[string]*Trigger),
	}
}

// Set sets a trigger, overwriting any existing trigger
func (t *Triggers) Set(name string, value any) {
	t.triggers[name] = NewTrigger(name, value)
}

// SetAfterSettle sets a trigger to be called after settle, overwriting any existing after-settle trigger
func (t *Triggers) SetAfterSettle(name string, value any) {
	t.afterSettle[name] = NewTrigger(name, value)
}

// SetAfterSwap sets a trigger to be called after swap, overwriting any existing after-swap trigger
func (t *Triggers) SetAfterSwap(name string, value any) {
	t.afterSwap[name] = NewTrigger(name, value)
}

// HasTriggers returns true if there are any triggers
func (t *Triggers) HasTriggers() bool {
	return len(t.triggers) > 0
}

// HasAfterSettleTriggers returns true if there are any after-settle triggers
func (t *Triggers) HasAfterSettleTriggers() bool {
	return len(t.afterSettle) > 0
}

// HasAfterSwapTriggers returns true if there are any after-swap triggers
func (t *Triggers) HasAfterSwapTriggers() bool {
	return len(t.afterSwap) > 0
}

// TriggerHeader returns all HTMX triggers as an HTTP header value
func (t *Triggers) TriggerHeader() (string, error) {
	return t.Encode(t.triggers)
}

// TriggerAfterSettleHeader returns the HTMX after-settle trigger as an HTTP header value
func (t *Triggers) TriggerAfterSettleHeader() (string, error) {
	return t.Encode(t.afterSettle)
}

// TriggerAfterSwapHeader returns the HTMX after-swap trigger as an HTTP header value
func (t *Triggers) TriggerAfterSwapHeader() (string, error) {
	return t.Encode(t.afterSwap)
}

// Encode encodes the triggers by marshalling them into an HTMX JSON object
func (t *Triggers) Encode(triggers map[string]*Trigger) (string, error) {
	events := make(map[string]any, len(triggers))

	for name, trigger := range triggers {
		if trigger == nil {
			continue
		}

		if trigger.value == nil {
			events[name] = ""
			continue
		}

		events[name] = trigger.value
	}

	if len(events) == 0 {
		return "", nil
	}

	bytes, err := json.Marshal(events)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
