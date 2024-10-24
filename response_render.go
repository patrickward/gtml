package gtml

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

// Render renders the response using the template manager
// Example: resp.StatusOK().Render(w, r)
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) {
	// Enforce a layout if none is set
	if resp.TemplateLayout() == "" {
		resp.Layout(resp.tm.baseLayout)
	}
	resp.tm.render(w, r, resp)
}

// RenderForbidden renders the 403 Forbidden page
func (resp *Response) RenderForbidden(w http.ResponseWriter, r *http.Request) {
	path := resp.tm.viewsPath(SystemDir, "403")
	if _, ok := resp.tm.templates[path]; ok {
		resp.Path(path).StatusForbidden().Render(w, r)
		return
	}
	resp.tm.handleError(w, r, ErrTempNotFound)
}

// RenderMaintenance renders the 503 Service Unavailable page
func (resp *Response) RenderMaintenance(w http.ResponseWriter, r *http.Request) {
	path := resp.tm.viewsPath(SystemDir, "503")
	if _, ok := resp.tm.templates[path]; ok {
		resp.Path(path).StatusUnavailable().Render(w, r)
		return
	}
	resp.tm.handleError(w, r, ErrTempNotFound)
}

// RenderMethodNotAllowed renders the 405 Method Not Allowed page
func (resp *Response) RenderMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	path := resp.tm.viewsPath(SystemDir, "405")
	if _, ok := resp.tm.templates[path]; ok {
		resp.Path(path).StatusError().Render(w, r)
		return
	}
	resp.tm.handleError(w, r, ErrTempNotFound)
}

// RenderNotFound renders the 404 Not Found page
func (resp *Response) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	path := resp.tm.viewsPath(SystemDir, "404")
	if _, ok := resp.tm.templates[path]; ok {
		resp.Path(path).StatusNotFound().Render(w, r)
		return
	}
	resp.tm.handleError(w, r, ErrTempNotFound)
}

// RenderSystemError renders the 500 Internal Server Error page
func (resp *Response) RenderSystemError(w http.ResponseWriter, r *http.Request, err error) {
	// Get the stack trace and output to the log
	resp.tm.logger.Error("Server error", slog.String("err", err.Error()))
	lineErrors := ""
	lines := strings.Split(string(debug.Stack()), "\n")
	for i, line := range lines {
		// replace \t with 4 spaces
		line = strings.ReplaceAll(line, "\t", "    ")
		lineErrors += fmt.Sprintf("--- traceLine%03d: %s\n", i, line)
		resp.tm.logger.Error("Stack trace", slog.String(fmt.Sprintf("--- traceLine%03d", i), line))
	}

	// If there is a template with the name "system/server_error" in the template cache, use it
	path := resp.tm.viewsPath(SystemDir, "500")
	if _, ok := resp.tm.templates[path]; ok {
		resp.Path(path).StatusError().Render(w, r)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// RenderUnauthorized renders the 401 Unauthorized page
func (resp *Response) RenderUnauthorized(w http.ResponseWriter, r *http.Request) {
	path := resp.tm.viewsPath(SystemDir, "401")
	if _, ok := resp.tm.templates[path]; ok {
		resp.Path(path).StatusUnauthorized().Render(w, r)
		return
	}
	resp.tm.handleError(w, r, ErrTempNotFound)
}
