package location

import (
	"encoding/json"
	"strings"

	"github.com/patrickward/gtml/htmx/swap"
)

type Option func(*Location)

// Event returns an Option that sets the event field of the Location struct
func Event(event string) Option {
	return func(l *Location) {
		l.Event = event
	}
}

// Handler returns an Option that sets the handler field of the Location struct
func Handler(handler string) Option {
	return func(l *Location) {
		l.Handler = handler
	}
}

// Headers returns an Option that sets the headers field of the Location struct
func Headers(headers map[string]string) Option {
	return func(l *Location) {
		l.Headers = headers
	}
}

// Select returns an Option that sets the select field of the Location struct
func Select(selectValue string) Option {
	return func(l *Location) {
		l.Select = selectValue
	}
}

// Source returns an Option that sets the source field of the Location struct
func Source(source string) Option {
	return func(l *Location) {
		l.Source = source
	}
}

// Swap returns an Option that sets the swap field of the Location struct
func Swap(swap *swap.Style) Option {
	return func(l *Location) {
		l.Swap = swap.String()
	}
}

// Target returns an Option that sets the target field of the Location struct
func Target(target string) Option {
	return func(l *Location) {
		l.Target = target
	}
}

// Values returns an Option that sets the values field of the Location struct
func Values(values map[string]string) Option {
	return func(l *Location) {
		l.Values = values
	}
}

// Location is a struct that represents a location to navigate to within HTMX.
// See the https://htmx.org/headers/hx-location/ documentation for more information.
type Location struct {
	// The path to navigate to
	Path string `json:"path"`

	// The event that triggered the request
	Event string `json:"event,omitempty"`

	// The type of request
	Handler string `json:"handler,omitempty"`

	// The headers to submit with the request
	Headers map[string]string `json:"headers,omitempty"`

	// The CSS selector to target the content to swap
	Select string `json:"select,omitempty"`

	// The source element of the request
	Source string `json:"source,omitempty"`

	// The type of swap to perform
	Swap string `json:"swap,omitempty"`

	// The target element to swap content into
	Target string `json:"target,omitempty"`

	// The values to submit with the request
	Values map[string]string `json:"values,omitempty"`

	pathOnly bool
}

// NewLocation creates a new Location struct with the given path and options
func NewLocation(path string, opt ...Option) *Location {
	l := new(Location)
	l.Path = path
	for _, o := range opt {
		o(l)
	}

	if opt == nil {
		l.pathOnly = true
	}
	return l
}

// String encodes the Location struct into a json string
func (hxl *Location) String() string {
	// If path is empty, then return an empty string
	if strings.TrimSpace(hxl.Path) == "" {
		return ""
	}

	// If the only value filled is path, then just return the path
	if hxl.pathOnly {
		return hxl.Path
	}

	data, err := json.Marshal(hxl)
	if err != nil {
		return ""
	}
	return string(data)
}
