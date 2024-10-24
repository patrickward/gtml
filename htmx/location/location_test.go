package location_test

import (
	"testing"

	"github.com/patrickward/gtml/htmx/location"
	"github.com/patrickward/gtml/htmx/swap"
)

func TestLocation_Encode(t *testing.T) {
	tests := []struct {
		name string
		hxl  *location.Location
		want string
	}{
		{
			name: "Path only",
			hxl:  location.NewLocation("/testpath"),
			want: "/testpath",
		},
		{
			name: "simple",
			// hxl:  location.Location{Path: "/testpath", Event: "click"},
			hxl:  location.NewLocation("/testpath", location.Event("click")),
			want: `{"path":"/testpath","event":"click"}`,
		},
		{
			name: "complex",
			hxl: location.NewLocation(
				"/testpath",
				location.Event("click"),
				location.Handler("handleClick"),
				location.Headers(map[string]string{"header": "value"}),
				location.Select("selectitem"),
				location.Source("source"),
				location.Swap(swap.InnerHTML(swap.Transition(true))),
				location.Target("target"),
				location.Values(map[string]string{"valueskey": "value"}),
			),
			want: `{"path":"/testpath","event":"click","handler":"handleClick","headers":{"header":"value"},"select":"selectitem","source":"source","swap":"innerHTML transition:true","target":"target","values":{"valueskey":"value"}}`,
		},
		{
			name: "empty",
			hxl:  location.NewLocation(""),
			want: "",
		},
		{
			name: "omit empty",
			hxl:  location.NewLocation("/pathonly", location.Event("eventonly")),
			want: `{"path":"/pathonly","event":"eventonly"}`,
		},
		{
			name: "only headers",
			hxl:  location.NewLocation("/", location.Headers(map[string]string{"header": "value"})),
			want: `{"path":"/","headers":{"header":"value"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hxl.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
