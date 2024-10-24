package gtml_test

import (
	"html/template"
	"io/fs"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patrickward/gtml"
	"github.com/patrickward/gtml/testdata/source1"
	"github.com/patrickward/gtml/testdata/source2"
)

type TestData struct {
	Title      string
	Content    string
	User       string
	IsAdmin    bool
	Navigation []string
}

func (td TestData) toMap() map[string]interface{} {
	return map[string]interface{}{
		"Title":      td.Title,
		"Content":    td.Content,
		"User":       td.User,
		"IsAdmin":    td.IsAdmin,
		"Navigation": td.Navigation,
	}
}

func TestTemplateManager(t *testing.T) {
	tests := []struct {
		name           string
		sources        map[string]fs.FS
		layout         string
		page           string
		data           TestData
		requestPath    string
		requestMethod  string
		requestHeaders map[string]string
		expectedParts  []string
		expectError    bool
	}{
		{
			name: "basic page with base layout from source1",
			sources: map[string]fs.FS{
				"": source1.FS,
			},
			layout: "base",
			page:   "home",
			data: TestData{
				Title:      "Welcome Home",
				Content:    "Main content here",
				User:       "John Doe",
				IsAdmin:    false,
				Navigation: []string{"Home", "About", "Contact"},
			},
			requestPath:   "/",
			requestMethod: "GET",
			requestHeaders: map[string]string{
				"Accept": "text/html",
			},
			expectedParts: []string{
				"<title>Welcome Home</title>",
				"Main content here",
				"John Doe",
				"<nav>Home</nav>",
			},
			expectError: false,
		},
		{
			name: "basic page with base layout from source1 and prefixed path",
			sources: map[string]fs.FS{
				"foobar": source1.FS,
			},
			layout: "base",
			page:   "foobar:home",
			data: TestData{
				Title:      "Welcome Home",
				Content:    "Main content here",
				User:       "John Doe",
				IsAdmin:    false,
				Navigation: []string{"Home", "About", "Contact"},
			},
			requestPath:   "/",
			requestMethod: "GET",
			requestHeaders: map[string]string{
				"Accept": "text/html",
			},
			expectedParts: []string{
				"<title>Welcome Home</title>",
				"Main content here",
				"John Doe",
				"<nav>Home</nav>",
			},
			expectError: false,
		},
		{
			name: "admin page with admin layout from source1",
			sources: map[string]fs.FS{
				"": source1.FS,
			},
			layout: "admin",
			page:   "admin/dash",
			data: TestData{
				Title:      "Admin Dashboard",
				Content:    "Dashboard Content",
				User:       "Admin User",
				IsAdmin:    true,
				Navigation: []string{"Dashboard", "Users", "Settings"},
			},
			requestPath:   "/admin/dashboard",
			requestMethod: "GET",
			requestHeaders: map[string]string{
				"Authorization": "Bearer test-token",
			},
			expectedParts: []string{
				"<title>Admin Dashboard</title>",
				"Dashboard Content",
				"Admin User",
				"<div class=\"admin-header\">",
			},
			expectError: false,
		},
		{
			name: "multiple sources with override",
			sources: map[string]fs.FS{
				"":        source1.FS,
				"source2": source2.FS,
			},
			layout: "source2:clean",
			page:   "home",
			data: TestData{
				Title:   "Clean Layout",
				Content: "Minimal content",
			},
			requestPath:   "/clean",
			requestMethod: "GET",
			requestHeaders: map[string]string{
				"Accept": "text/html",
			},
			expectedParts: []string{
				"<main class=\"clean-layout-source2\">",
				"Minimal content",
				"Source 2 Header",
			},
			expectError: false,
		},
		{
			name: "non-existent layout",
			sources: map[string]fs.FS{
				"": source1.FS,
			},
			layout:         "missing",
			page:           "home",
			data:           TestData{},
			requestPath:    "/",
			requestMethod:  "GET",
			requestHeaders: map[string]string{},
			expectedParts: []string{
				"error executing template: html/template: \"layout:missing\"",
			},
			expectError: false,
		},
		{
			name: "non-existent page",
			sources: map[string]fs.FS{
				"": source1.FS,
			},
			layout:         "base",
			page:           "missing",
			data:           TestData{},
			requestPath:    "/missing",
			requestMethod:  "GET",
			requestHeaders: map[string]string{},
			expectedParts: []string{
				"template not found",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		logger := slog.New(slog.NewTextHandler(httptest.NewRecorder(), nil))
		funcMap := template.FuncMap{
			"isAdmin": func(user string) bool {
				return strings.HasPrefix(user, "Admin")
			},
		}
		t.Run(tt.name, func(t *testing.T) {
			// Initialize template manager and load templates
			tm, err := gtml.NewTemplateManager(gtml.TemplateManagerOptions{
				Extension: ".gtml",
				Sources:   tt.sources,
				Funcs:     funcMap,
				Logger:    logger,
			})

			require.NoError(t, err, "Failed to load templates")

			// Test rendering
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, nil)
			for key, value := range tt.requestHeaders {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()
			tm.NewResponse().
				Layout(tt.layout).
				Path(tt.page).
				Data(tt.data.toMap()).
				Title(tt.data.Title).
				Render(w, req)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			result := w.Body.String()

			for _, expected := range tt.expectedParts {
				assert.Contains(t, result, expected,
					"Expected content not found in rendered template")
			}

			// TODO: test headers?
		})
	}
}
