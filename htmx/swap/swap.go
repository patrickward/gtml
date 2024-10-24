package swap

import (
	"fmt"
	"time"
)

type Option func(*Style)

// Direction represents the possible scroll direction values for the show and scroll modifiers
type Direction string

const (
	DirectionTop    Direction = "top"
	DirectionBottom Direction = "bottom"
)

// Style represents an HTMX swap style that can be used to instruct HTMX how to swap content.
//
// For more information, see: https://htmx.org/attributes/hx-swap
type Style struct {
	style       string
	transition  string
	swap        string
	settle      string
	ignoreTitle string
	scroll      string
	show        string
	focusScroll string
}

func newStyle(style string, opt ...Option) *Style {
	s := &Style{}
	s.style = style
	for _, o := range opt {
		o(s)
	}
	return s
}

// InnerHTML replaces the inner HTML of the target element
func InnerHTML(opt ...Option) *Style {
	return newStyle("innerHTML", opt...)
}

// OuterHTML replaces the entire target element with the response
func OuterHTML(opt ...Option) *Style {
	return newStyle("outerHTML", opt...)
}

// BeforeBegin Inserts the response before the target element
func BeforeBegin(opt ...Option) *Style {
	return newStyle("beforebegin", opt...)
}

// AfterBegin Inserts the response before the first child of the target element
func AfterBegin(opt ...Option) *Style {
	return newStyle("afterbegin", opt...)
}

// BeforeEnd Inserts the response after the last child of the target element
func BeforeEnd(opt ...Option) *Style {
	return newStyle("beforeend", opt...)
}

// AfterEnd Inserts the response after the target element
func AfterEnd(opt ...Option) *Style {
	return newStyle("afterend", opt...)
}

// Delete Deletes the target element regardless of the response
func Delete(opt ...Option) *Style {
	return newStyle("delete", opt...)
}

// None Instructs HTMX not to append content from the response. However, out of band swaps will still be processed.
func None(opt ...Option) *Style {
	return newStyle("none", opt...)
}

// String returns the string representation of the Style by combining the style and space separated modifiers
func (s *Style) String() string {
	output := s.style
	output = s.append(output, "transition", s.transition)
	output = s.append(output, "swap", s.swap)
	output = s.append(output, "settle", s.settle)
	output = s.append(output, "ignoreTitle", s.ignoreTitle)
	output = s.append(output, "scroll", s.scroll)
	output = s.append(output, "show", s.show)
	output = s.append(output, "focus-scroll", s.focusScroll)
	return output
}

func Transition(val bool) Option {
	return func(s *Style) {
		s.transition = boolToStr(val)
	}
}

func IgnoreTitle() Option {
	return func(s *Style) {
		s.ignoreTitle = "true"
	}
}

//goland:noinspection GoNameStartsWithPackageName
func SwapAfter(d time.Duration) Option {
	return func(s *Style) {
		s.swap = d.String()
	}
}

func SettleAfter(d time.Duration) Option {
	return func(s *Style) {
		s.settle = d.String()
	}
}

func Scroll(direction Direction) Option {
	return func(s *Style) {
		s.scroll = string(direction)
	}
}

func ScrollTo(selector string, direction Direction) Option {
	return func(s *Style) {
		s.scroll = fmt.Sprintf("%s:%s", selector, string(direction))
	}
}

func Show(direction Direction) Option {
	return func(s *Style) {
		s.show = string(direction)
	}
}

func ShowTo(selector string, direction Direction) Option {
	return func(s *Style) {
		s.show = fmt.Sprintf("%s:%s", selector, string(direction))
	}
}

func ShowNone() Option {
	return func(s *Style) {
		s.show = "none"
	}
}

func FocusScroll(val bool) Option {
	return func(s *Style) {
		s.focusScroll = boolToStr(val)
	}
}

func boolToStr(val bool) string {
	if val {
		return "true"
	}
	return "false"
}

// append checks if field is not empty and appends it to value
func (s *Style) append(output, fieldName, value string) string {
	if value != "" {
		output += fmt.Sprintf(" %s:%s", fieldName, value)
	}
	return output
}
