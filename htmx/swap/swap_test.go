package swap_test

import (
	"testing"
	"time"

	"github.com/patrickward/gtml/htmx/swap"
)

type swapTest struct {
	name     string
	swapFunc func() *swap.Style
	// swapOptions func(*swap.Style) *swap.Style
	expected string
}

func TestSwapOption_Transition(t *testing.T) {
	tests := []swapTest{
		{
			name: "WithTransition true",
			swapFunc: func() *swap.Style {
				return swap.OuterHTML(swap.Transition(true))
			},
			expected: "outerHTML transition:true",
		},
		{
			name: "WithTransition false",
			swapFunc: func() *swap.Style {
				return swap.OuterHTML(swap.Transition(false))
			},
			expected: "outerHTML transition:false",
		},
		{
			name: "WithIgnoreTitle",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.IgnoreTitle())
			},
			expected: "innerHTML ignoreTitle:true",
		},
	}

	runSwapTests(t, tests)
}

func TestSwapOption_Show(t *testing.T) {
	tests := []swapTest{
		{
			name: "WithShowTop",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.Show(swap.DirectionTop))
			},
			expected: "innerHTML show:top",
		},
		{
			name: "WithShowToBottom",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.ShowTo("#foobar", swap.DirectionBottom))
			},
			expected: "innerHTML show:#foobar:bottom",
		},
		{
			name: "WithShowNone",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.ShowNone())
			},
			expected: "innerHTML show:none",
		},
	}

	runSwapTests(t, tests)
}

func TestSwapOption_Scroll(t *testing.T) {
	tests := []swapTest{
		{
			name: "WithScrollTop",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.Scroll(swap.DirectionTop))
			},
			expected: "innerHTML scroll:top",
		},
		{
			name: "WithScrollToBottom",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.ScrollTo("#foobar", swap.DirectionBottom))
			},
			expected: "innerHTML scroll:#foobar:bottom",
		},
	}

	runSwapTests(t, tests)
}

func TestSwapOption_IgnoreTitle(t *testing.T) {
	tests := []swapTest{
		{
			name: "WithIgnoreTitle",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.IgnoreTitle())
			},
			expected: "innerHTML ignoreTitle:true",
		},
	}

	runSwapTests(t, tests)
}

func TestSwapOption_SwapAfter(t *testing.T) {
	tests := []swapTest{
		{
			name: "WithSwapAfter",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.SwapAfter(time.Second * 2))
			},
			expected: "innerHTML swap:2s",
		},
	}

	runSwapTests(t, tests)
}

func TestSwapOption_SettleAfter(t *testing.T) {
	tests := []swapTest{
		{
			name: "WithSettleAfter",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.SettleAfter(time.Second * 4))
			},
			expected: "innerHTML settle:4s",
		},
	}

	runSwapTests(t, tests)
}

func TestSwapOption_FocusScroll(t *testing.T) {
	tests := []swapTest{
		{
			name: "WithFocusScroll",
			swapFunc: func() *swap.Style {
				return swap.InnerHTML(swap.FocusScroll(true))
			},
			expected: "innerHTML focus-scroll:true",
		},
	}

	runSwapTests(t, tests)
}

func TestSwapOption_MultipleSwapModifiers(t *testing.T) {
	tests := []swapTest{
		{
			name: "WithMultipleSwapModifiers",
			swapFunc: func() *swap.Style {
				return swap.OuterHTML(
					swap.SettleAfter(time.Second*4),
					swap.Scroll(swap.DirectionBottom),
					swap.FocusScroll(true),
				)
			},
			expected: "outerHTML settle:4s scroll:bottom focus-scroll:true",
		},
		{
			name: "WithAllSwapModifiers",
			swapFunc: func() *swap.Style {
				return swap.BeforeEnd(
					swap.Transition(true),
					swap.SwapAfter(time.Second*2),
					swap.SettleAfter(time.Second*4),
					swap.IgnoreTitle(),
					swap.ScrollTo("#foobarbaz", swap.DirectionBottom),
					swap.ShowTo("#foobarbaz", swap.DirectionTop),
					swap.FocusScroll(false),
				)
			},
			expected: "beforeend transition:true swap:2s settle:4s ignoreTitle:true scroll:#foobarbaz:bottom show:#foobarbaz:top focus-scroll:false",
		},
	}
	runSwapTests(t, tests)
}

func runSwapTests(t *testing.T, tests []swapTest) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.swapFunc()
			got := s.String()
			if got != tt.expected {
				t.Errorf("Expected \"%v\", got %v", tt.expected, got)
			}
		})
	}
}
