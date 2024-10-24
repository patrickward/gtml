package gtml

import (
	"context"
	"net/http"
	"time"

	"github.com/patrickward/gtml/request"
)

const NonceContextKey = "gtml_nonce"

// ResponseData is the struct that all view models must implement. It provides common data for all templates
// and represents the data that is passed to the template.
//
// This is a short-lived object that is used to work with data passed to the template. It is not thread-safe.
//
// ResponseData should not be used directly. Instead, use the NewResponseData function to create an instance
// of ResponseData that contains the data you want to pass to the template.
//
// Example: NewResponseData(request, map[string]any{"title": "Hello World"})
//
// The environment variables are generally added via config/config, but if you're not using that package,
// you can set them manually.
//
//goland:noinspection GoNameStartsWithPackageName
type ResponseData struct {
	title       string
	request     *http.Request
	pageData    map[string]any
	csrfToken   string
	environment string
}

// NewResponseData creates a new Data instance.
// If you are using this outside the normal HyperView rendering process, be sure to set the request manually
// via Data.SetRequest as the request is deliberately set later in the normal rendering flow.
func NewResponseData(pageData map[string]any) *ResponseData {
	pageData = initData(pageData)
	return &ResponseData{
		pageData: pageData,
	}
}

// SetTitle sets the title of the page.
func (v *ResponseData) SetTitle(title string) {
	v.title = title
}

// SetRequest sets the request for the Data instance.
func (v *ResponseData) SetRequest(r *http.Request) {
	v.request = r
}

func initData(data map[string]any) map[string]any {
	if data == nil {
		data = map[string]any{}
	}

	// If no "Error" key is set, set it to an empty string
	if _, ok := data["Error"]; !ok {
		data["Error"] = ""
	}

	// If no "Errors" key is set, set it to an empty map
	if _, ok := data["Errors"]; !ok {
		data["Errors"] = map[string]string{}
	}

	return data
}

// Data returns the data map that will be passed to the template.
func (v *ResponseData) Data() map[string]any {
	v.pageData = initData(v.pageData)
	v.pageData["View"] = v
	return v.pageData
}

// AddData adds a map of data to the existing view data model.
func (v *ResponseData) AddData(data map[string]any) {
	for key, value := range data {
		v.pageData[key] = value
	}
}

// AddDataItem adds a single key-value pair to the existing view data model.
func (v *ResponseData) AddDataItem(key string, value any) {
	v.pageData[key] = value
}

// AddErrors adds an error message and a map of field errors to the view data model.
func (v *ResponseData) AddErrors(msg string, fieldErrors map[string]string) {
	v.pageData["Error"] = msg
	v.pageData["Errors"] = fieldErrors
}

// Get returns the value of the specified key from the view data model.
func (v *ResponseData) Get(key string) any {
	val, ok := v.pageData[key]
	if ok {
		return val
	}

	return ""
}

// GetString returns the value of the specified key from the view data model as a string.
func (v *ResponseData) GetString(key string) string {
	val, ok := v.Get(key).(string)
	if ok {
		return val
	}

	return ""
}

// Title returns the title of the page.
func (v *ResponseData) Title() string {
	return v.title
}

// ------ Error Helpers --------

// HasError returns true if the view data model contains an error message.
func (v *ResponseData) HasError() bool {
	return v.GetString("Error") != ""
}

// Error returns the error message from the view data model.
func (v *ResponseData) Error() string {
	return v.GetString("Error")
}

// HasErrors returns true if the view data model contains field errors.
func (v *ResponseData) HasErrors() bool {
	return len(v.Errors()) > 0
}

// Errors returns a map of field errors from the view data model.
func (v *ResponseData) Errors() map[string]string {
	val, ok := v.Get("Errors").(map[string]string)
	if ok {
		return val
	}

	return map[string]string{}
}

// ------ Common Helpers --------

// BaseURL returns the base URL of the request.
func (v *ResponseData) BaseURL() string {
	return request.BaseURL(v.request)
}

// Context returns the context of the request.
func (v *ResponseData) Context() context.Context {
	return v.request.Context()
}

// CurrentYear returns the current year.
func (v *ResponseData) CurrentYear() int {
	return time.Now().Year()
}

// Nonce returns the nonce value from the request context, if available.
func (v *ResponseData) Nonce() string {
	nonce, ok := v.request.Context().Value(NonceContextKey).(string)
	if ok {
		return nonce
	}

	return ""
}

// RequestPath returns the path of the request.
func (v *ResponseData) RequestPath() string {
	return request.URLPath(v.request)
}

// RequestMethod returns the method of the request.
func (v *ResponseData) RequestMethod() string {
	return request.Method(v.request)
}
