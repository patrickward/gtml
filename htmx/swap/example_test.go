package swap_test

import (
	"fmt"
	"time"

	"github.com/patrickward/gtml/htmx/swap"
)

func ExampleInnerHTML() {
	// For simple swaps, use any of the style functions alone
	s1 := swap.InnerHTML()

	// For more complex swaps, use any of the style functions with the options pattern
	s2 := swap.AfterBegin(
		swap.Transition(true),
		swap.IgnoreTitle(),
		swap.SwapAfter(time.Second*2),
	)

	fmt.Println(s1.String())
	fmt.Println(s2.String())

	// Output:
	// innerHTML
	// afterbegin transition:true swap:2s ignoreTitle:true
}
