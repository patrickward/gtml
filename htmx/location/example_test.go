package location_test

import (
	"fmt"

	"github.com/patrickward/gtml/htmx/location"
	"github.com/patrickward/gtml/htmx/swap"
)

func ExampleNewLocation() {
	// For simple location strings, use NewLocation with a single path parameter
	loc1 := location.NewLocation("/testpath")

	// For more complex location strings with context, use NewLocation with the options pattern
	loc2 := location.NewLocation(
		"/testpath",
		location.Event("click"),
		location.Handler("handleClick"),
		location.Headers(map[string]string{"header": "value"}),
		location.Select("#foobar"),
		location.Source("source"),
		location.Swap(swap.InnerHTML(swap.Transition(true), swap.IgnoreTitle())),
		location.Target("target"),
		location.Values(map[string]string{"foo": "bar"}),
	)

	fmt.Println(loc1.String())
	fmt.Println(loc2.String())

	// Output:
	// /testpath
	// {"path":"/testpath","event":"click","handler":"handleClick","headers":{"header":"value"},"select":"#foobar","source":"source","swap":"innerHTML transition:true ignoreTitle:true","target":"target","values":{"foo":"bar"}}
}
