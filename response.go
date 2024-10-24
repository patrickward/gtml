package gtml

import (
	"net/http"
	"strings"

	"github.com/patrickward/gtml/htmx"
	"github.com/patrickward/gtml/htmx/trigger"
)

// Response represents a view response to an HTTP request
// It uses a fluent interface to allow for chaining of methods, so that methods can be called in any order.
type Response struct {
	// The headers to be passed to the response (default: empty)
	headers map[string]string
	// The layout template to be used (required, no default)
	layout string
	// The view template path to be used (required, no default)
	path string
	// The status code to be passed to the response (default: http.StatusOK)
	request *http.Request
	// The status code to be passed to the response (default: http.StatusOK)
	statusCode int
	// The title of the page (default: the page name without the extension)
	title string
	// The triggers to be passed to the response (default: empty)
	triggers *trigger.Triggers
	// The view data to be passed to the template (default: ViewData{})
	data *ResponseData
	// The template manager to be used for rendering templates
	tm *TemplateManager
}

func NewResponse(tm *TemplateManager) *Response {
	return &Response{
		data:       NewResponseData(make(map[string]any)),
		headers:    map[string]string{},
		layout:     "",
		path:       "",
		request:    nil,
		statusCode: http.StatusOK,
		title:      "",
		triggers:   trigger.NewTriggers(),
		tm:         tm,
	}
}

// ViewData returns the view data model. The request is set here to ensure
// the request is available in the template and that it is not overwritten until later in the process.
func (resp *Response) ViewData(r *http.Request) *ResponseData {
	resp.data.SetTitle(resp.title)
	resp.data.SetRequest(r)
	return resp.data
}

// Headers returns the headers map as a combination map of both triggers and headers
func (resp *Response) Headers() map[string]string {
	if resp.headers == nil {
		resp.headers = map[string]string{}
	}

	if resp.triggers != nil {
		if resp.triggers.HasTriggers() {
			val, err := resp.triggers.TriggerHeader()
			if err == nil {
				resp.headers[htmx.HXTrigger] = val
			}
		}

		if resp.triggers.HasAfterSettleTriggers() {
			val, err := resp.triggers.TriggerAfterSwapHeader()
			if err == nil {
				resp.headers[htmx.HXTriggerAfterSwap] = val
			}
		}

		if resp.triggers.HasAfterSwapTriggers() {
			val, err := resp.triggers.TriggerAfterSettleHeader()
			if err == nil {
				resp.headers[htmx.HXTriggerAfterSettle] = val
			}
		}
	}

	return resp.headers
}

// HTTPHeader returns a http.Header for the headers map
func (resp *Response) HTTPHeader() http.Header {
	if resp.headers == nil {
		return nil
	}

	header := make(http.Header)
	for key, value := range resp.Headers() {
		header.Set(key, value)
	}

	return header
}

// TemplateLayout returns the template layout
func (resp *Response) TemplateLayout() string {
	return resp.layout
}

// TemplatePath returns the path used in templates, if any
func (resp *Response) TemplatePath() string {
	return resp.path
}

// PageTitle returns the page title
func (resp *Response) PageTitle() string {
	return resp.title
}

// StatusCode returns the status code.
func (resp *Response) StatusCode() int {
	return resp.statusCode
}

// ResetData resets the view data model with an existing model. It returns the modified Response pointer.
// The view data model contains data that will be passed to the view template for rendering.
//
// Alternatively, you can add the data map and create a new view data model automatically using the Data function.
//
// Important: if you are creating the view data model externally and need to use it before render is called,
// you should probably set the request via Data.SetRequest, as the request is deliberately set later in
// the rendering process in most cases.
func (resp *Response) ResetData(data *ResponseData) *Response {
	resp.data = data
	return resp
}

// Data creates a new view data model with the provided data map and returns the modified Response pointer.
// This will overwrite any existing view data model. If you want to add data to an existing view data model, create
// a new view data model externally using the NewResponseData function and pass it to the ResetData function instead.
func (resp *Response) Data(data map[string]any) *Response {
	resp.data = NewResponseData(data)
	return resp
}

// AddData adds data to the view data model. It returns the modified Response pointer.
func (resp *Response) AddData(data map[string]any) *Response {
	resp.data.AddData(data)
	return resp
}

// AddDataItem adds a single data item to the view data model. It returns the modified Response pointer.
func (resp *Response) AddDataItem(key string, value any) *Response {
	resp.data.AddDataItem(key, value)
	return resp
}

// Errors adds an error message and any field errors to the view data model.
// This will also set the status code to 422 (Unprocessable Entity)). If that is not correct status code,
// you should reset it using the Status() function or one of the Status* shortcut functions.
func (resp *Response) Errors(msg string, fieldErrors map[string]string) *Response {
	resp.statusCode = http.StatusUnprocessableEntity
	resp.data.AddErrors(msg, fieldErrors)
	return resp
}

// Title sets the page title
func (resp *Response) Title(title string) *Response {
	resp.title = title
	return resp
}

// Path sets the template path
func (resp *Response) Path(path string) *Response {
	// If the path contains a colon, it's part of a plugin path, so we need to
	// extract the plugin name from the path first
	pathParts := strings.SplitN(path, ":", 2)

	if len(pathParts) == 2 {
		path = pathParts[1]
	}

	if !strings.HasPrefix(path, ViewsDir+"/") {
		path = ViewsDir + "/" + path
	}

	if len(pathParts) == 2 {
		path = pathParts[0] + ":" + path
	}

	resp.path = path
	return resp
}

// Layout sets the template layout. It updates the layout value in the Response struct.
// Then it returns the updated Response struct itself for method chaining.
func (resp *Response) Layout(layout string) *Response {
	resp.layout = layout
	return resp
}

// Header adds/sets a header
func (resp *Response) Header(key, value string) *Response {
	if resp.headers == nil {
		resp.headers = make(map[string]string)
	}

	resp.headers[key] = value
	return resp
}

// Status sets the status code.
func (resp *Response) Status(status int) *Response {
	resp.statusCode = status
	return resp
}

// NoCacheStrict sets the Cache-Control header to "no-cache, no-store, must-revalidate".
func (resp *Response) NoCacheStrict() {
	resp.headers["Cache-Control"] = "no-cache, no-store, must-revalidate"
}

// CacheControl sets the Cache-Control header to the given value.
func (resp *Response) CacheControl(cacheControl string) {
	resp.headers["Cache-Control"] = cacheControl
}

// ETag sets the ETag header to the given value.
func (resp *Response) ETag(etag string) {
	resp.headers["ETag"] = etag
}

// LastModified sets the Last-Modified header to the given value.
func (resp *Response) LastModified(lastModified string) {
	resp.headers["Last-Modified"] = lastModified
}
