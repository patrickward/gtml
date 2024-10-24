package gtml

import (
	"fmt"
	"html/template"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// TODO: Add test cases for all functions

const (
	day  = 24 * time.Hour
	year = 365 * day
)

var caser = cases.Title(language.English)
var printer = message.NewPrinter(language.English)

// MergeFuncMaps merges the provided function map with the default function map.
// If there are any conflicts, the provided function map will take precedence.
func MergeFuncMaps(funcMap template.FuncMap) template.FuncMap {
	defaultFuncMap := DefaultFuncMap()
	for key, value := range funcMap {
		defaultFuncMap[key] = value
	}
	return defaultFuncMap
}

func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		// Maps/Collections
		"kv":    kv,           // Create a ke-value pair
		"dict":  dict,         // Create a map from key-value pairs
		"slice": slice,        // Create a slice from a list of items
		"group": group,        // Split a slice into chunks of specified size
		"first": first,        // Get the first element of a slice
		"last":  last,         // Get the last element of a slice
		"nth":   nth,          // Get the nth element of a slice
		"join":  strings.Join, // Join a slice into a string

		// TODO: Implement the following functions
		//"shuffle":   shuffle,      // Shuffle a slice
		//"sort":      sort,         // Sort a slice
		//"reverse":   reverse,      // Reverse a slice
		//"unique":    unique,       // Get unique elements from a slice
		//"pluck":     pluck,        // Get a slice of values from a map
		//"hasString": hasString,    // Check if a string exists in a slice
		//"hasInt":    hasInt,       // Check if an int exists in a slice

		// String operations
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"trim":      strings.Trim,
		"trimSpace": strings.TrimSpace,
		"split":     strings.Split,
		"contains":  strings.Contains,
		"titleize":  titleize,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,

		// Time functions
		"now":       time.Now,
		"fmtTime":   formatTime,
		"isToday":   isToday,
		"timeSince": time.Since,
		"timeUntil": time.Until,
		"duration":  approximateDuration,

		// Type conversions
		"toString": toString,
		"toInt":    toInt,
		"toFloat":  toFloat,

		// Number operations
		"fmtInt":   formatInt,
		"fmtFloat": formatFloat,

		// Math/Logic
		"incr": incr, // Increment a number by 1
		"decr": decr, // Decrement a number by 1

		// Boolean functions
		"yesno": yesno,

		// HTML functions
		"safe":        safeHTML,
		"urlSetParam": urlSetParam,
		"urlDelParam": urlDelParam,
	}
}

// kv is a helper function that creates a key-value pair and merges it with another map, if provided.
func kv(k string, v any, other map[string]any) map[string]any {
	result := map[string]any{k: v}
	if other != nil {
		for key, value := range other {
			result[key] = value
		}
	}
	return result
}

// dict creates a map from a list of key-value pairs.
func dict(pairs ...any) (map[string]any, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("invalid number of arguments, there must be an even number of arguments")
	}
	result := make(map[string]any)
	for i := 0; i < len(pairs); i += 2 {
		if i+1 < len(pairs) {
			key := pairs[i].(string)
			value := pairs[i+1]
			result[key] = value
		}
	}
	return result, nil
}

// slice creates a slice from a list of items.
func slice(items ...any) []any {
	return items
}

// group splits a slice into chunks of specified size
func group(size int, seq []interface{}) [][]interface{} {
	if size <= 0 {
		return nil
	}
	chunks := make([][]interface{}, 0, (len(seq)+size-1)/size)
	for size < len(seq) {
		seq, chunks = seq[size:], append(chunks, seq[0:size:size])
	}
	chunks = append(chunks, seq)
	return chunks
}

// first returns the first element of a slice
func first(seq []any) any {
	if len(seq) > 0 {
		return seq[0]
	}
	return nil
}

// last returns the last element of a slice
func last(seq []any) any {
	if len(seq) > 0 {
		return seq[len(seq)-1]
	}
	return nil
}

// nth returns the nth element of a slice
func nth(seq []any, n int) any {
	if n < 0 || n >= len(seq) {
		return nil
	}
	return seq[n]
}

// toString converts any type to a string
func toString(i any) string {
	return fmt.Sprintf("%v", i)
}

// titleize capitalizes the first letter of each word in a string (English)
func titleize(s string) string {
	return caser.String(s)
}

// toInt converts any type to an int
func toInt(i any) (int, error) {
	switch v := i.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", i)
	}
}

// toFloat converts any type to a float64
func toFloat(i any) (float64, error) {
	switch v := i.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", i)
	}
}

// Increment an integer value
func incr(i any) (int64, error) {
	n, err := toInt64(i)
	if err != nil {
		return 0, err
	}

	n++
	return n, nil
}

// Decrement an integer value
func decr(i any) (int64, error) {
	n, err := toInt64(i)
	if err != nil {
		return 0, err
	}

	n--
	return n, nil
}

// toInt64 converts any type to int64
func toInt64(i any) (int64, error) {
	switch v := i.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	// Note: uint64 not supported due to risk of truncation.
	case string:
		return strconv.ParseInt(v, 10, 64)
	}

	return 0, fmt.Errorf("unable to convert type %T to int", i)
}

// isToday checks if the given time is today
func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}

// formatInt formats an integer as a string
func formatInt(i any) (string, error) {
	n, err := toInt64(i)
	if err != nil {
		return "", err
	}

	return printer.Sprintf("%d", n), nil
}

// formatFloat formats a float64 as a string with specified decimal places
func formatFloat(f float64, dp int) string {
	format := "%." + strconv.Itoa(dp) + "f"
	return printer.Sprintf(format, f)
}

// safeHTML returns a template.HTML value to mark a string as safe for HTML output
func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

// yesno returns "Yes" or "No" based on the boolean value
func yesno(b bool) string {
	if b {
		return "Yes"
	}

	return "No"
}

// formatTime formats a time.Time value as a string
func formatTime(format string, t time.Time) string {
	return t.Format(format)
}

// urlSetParam sets a query parameter in a URL
func urlSetParam(u *url.URL, key string, value any) *url.URL {
	nu := *u
	values := nu.Query()

	values.Set(key, fmt.Sprintf("%v", value))

	nu.RawQuery = values.Encode()
	return &nu
}

// urlDelParam deletes a query parameter from a URL
func urlDelParam(u *url.URL, key string) *url.URL {
	nu := *u
	values := nu.Query()

	values.Del(key)

	nu.RawQuery = values.Encode()
	return &nu
}

// approximateDuration returns a human-readable approximation of a duration
func approximateDuration(d time.Duration) string {
	if d < time.Second {
		return "less than 1 second"
	}

	ds := int(math.Round(d.Seconds()))
	if ds == 1 {
		return "1 second"
	} else if ds < 60 {
		return fmt.Sprintf("%d seconds", ds)
	}

	dm := int(math.Round(d.Minutes()))
	if dm == 1 {
		return "1 minute"
	} else if dm < 60 {
		return fmt.Sprintf("%d minutes", dm)
	}

	dh := int(math.Round(d.Hours()))
	if dh == 1 {
		return "1 hour"
	} else if dh < 24 {
		return fmt.Sprintf("%d hours", dh)
	}

	dd := int(math.Round(float64(d / day)))
	if dd == 1 {
		return "1 day"
	} else if dd < 365 {
		return fmt.Sprintf("%d days", dd)
	}

	dy := int(math.Round(float64(d / year)))
	if dy == 1 {
		return "1 year"
	}

	return fmt.Sprintf("%d years", dy)
}
