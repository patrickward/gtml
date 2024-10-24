package gtml

import "net/http"

// StatusCreated sets the status code to Created (201)
func (resp *Response) StatusCreated() *Response {
	resp.statusCode = http.StatusCreated
	return resp
}

// StatusAccepted sets the status code to Accepted (202) and returns the Response object.
func (resp *Response) StatusAccepted() *Response {
	resp.statusCode = http.StatusAccepted
	return resp
}

// StatusNoContent sets the status code to NoContent (204). It returns the Response object for method chaining.
func (resp *Response) StatusNoContent() *Response {
	resp.statusCode = http.StatusNoContent
	return resp
}

// StatusNotFound sets the status code to NotFound (404)
func (resp *Response) StatusNotFound() *Response {
	resp.statusCode = http.StatusNotFound
	return resp
}

// StatusForbidden sets the status code to Forbidden (403) and returns the Response object.
func (resp *Response) StatusForbidden() *Response {
	resp.statusCode = http.StatusForbidden
	return resp
}

// StatusUnavailable sets the status code to ServiceUnavailable (503)
func (resp *Response) StatusUnavailable() *Response {
	resp.statusCode = http.StatusServiceUnavailable
	return resp
}

// StatusUnprocessable sets the status code to UnprocessableEntity (422)
func (resp *Response) StatusUnprocessable() *Response {
	resp.statusCode = http.StatusUnprocessableEntity
	return resp
}

// StatusError sets the status code to InternalServerError (500)
func (resp *Response) StatusError() *Response {
	resp.statusCode = http.StatusInternalServerError
	return resp
}

// StatusOK sets the status code to OK (200)
func (resp *Response) StatusOK() *Response {
	resp.statusCode = http.StatusOK
	return resp
}

// StatusUnauthorized sets the status code to Unauthorized (401)
func (resp *Response) StatusUnauthorized() *Response {
	resp.statusCode = http.StatusUnauthorized
	return resp
}

// StatusStopPolling sets the status code to 286 and returns the Response object.
// This is useful when working HTMX and polling. Responding with a status of 286 will tell HTMX to stop polling.
// SEE: https://htmx.org/docs/#polling
func (resp *Response) StatusStopPolling() *Response {
	resp.statusCode = 286
	return resp
}
