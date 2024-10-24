package gtml

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
)

// TemplateManager is a template adapter for the HyperView framework that uses the Go html/template package.
type TemplateManager struct {
	baseLayout    string
	systemLayout  string
	extension     string
	fileSystemMap map[string]fs.FS
	logger        *slog.Logger
	funcMap       template.FuncMap
	templates     map[string]*template.Template
}

// TemplateManagerOptions are the options for the TemplateManager.
type TemplateManagerOptions struct {
	// BaseLayout is the default layout to use for rendering templates. Default is "base".
	BaseLayout string

	// SystemLayout is the layout to use for system pages (e.g. 404, 500). Default is "base".
	SystemLayout string

	// Extension is the file extension for the templates. Default is ".html".
	Extension string

	// Sources is a map of file systems to use for the templates. The string key is also used as a prefix for the template names.
	Sources map[string]fs.FS

	// Funcs is a map of functions to add to the template.FuncMap.
	Funcs template.FuncMap

	// Logger is the logger to use for the adapter.
	Logger *slog.Logger
}

// NewTemplateManager creates a new TemplateManager.
func NewTemplateManager(opts TemplateManagerOptions) *TemplateManager {
	funcMap := MergeFuncMaps(opts.Funcs)

	// Set default extension if not provided
	if opts.Extension == "" {
		opts.Extension = ".html"
	}

	// Ensure the extension starts with a .
	if opts.Extension[0] != '.' {
		opts.Extension = "." + opts.Extension
	}

	// If no base layout is provided, set it to "base"
	if opts.BaseLayout == "" {
		opts.BaseLayout = DefaultBaseLayout
	}

	// If no system layout is provided, set it to "base"
	if opts.SystemLayout == "" {
		opts.SystemLayout = opts.BaseLayout
	}

	return &TemplateManager{
		baseLayout:    opts.BaseLayout,
		systemLayout:  opts.SystemLayout,
		extension:     opts.Extension,
		fileSystemMap: opts.Sources,
		funcMap:       funcMap,
		logger:        opts.Logger,
		templates:     make(map[string]*template.Template),
	}
}

// NewResponse creates a new Response instance with the TemplateManager.
func (tm *TemplateManager) NewResponse() *Response {
	return NewResponse(tm)
}

func (tm *TemplateManager) Init() error {
	// Reset the template cache
	tm.templates = make(map[string]*template.Template)

	layoutsAndPartials, err := tm.loadLayoutsAndPartials()
	if err != nil {
		return fmt.Errorf("error loading partials. %w", err)
	}

	// Recursively process directories from all Sources
	for fsID, fsys := range tm.fileSystemMap {
		processDirectory := func(path string, dir fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !dir.IsDir() && filepath.Ext(path) == tm.extension {
				relPath, err := filepath.Rel("", path)
				if err != nil {
					return err
				}
				pageName := strings.TrimSuffix(relPath, filepath.Ext(relPath))
				//if fsID != RootFSID {
				if fsID != "" && fsID != "-" {
					pageName = fsID + ":" + pageName
				}

				// Clone the layout and partial templates and parse the page template,
				// so we can reuse the common templates for variants
				tmpl, err := template.Must(layoutsAndPartials.Clone()).ParseFS(fsys, path)

				if err != nil {
					return err
				}

				tm.templates[pageName] = tmpl
			}
			return nil
		}

		// If the "views" directory exists, parse it.
		if _, err := fsys.Open(ViewsDir); err == nil {
			if err := fs.WalkDir(fsys, ViewsDir, processDirectory); err != nil {
				return err
			}
		}
	}

	// Uncomment to view the template names found
	//tm.printTemplateNames()

	return nil
}

func (tm *TemplateManager) loadLayoutsAndPartials() (*template.Template, error) {
	commonTemplates := template.New("_common_").Funcs(tm.funcMap)

	for _, fsys := range tm.fileSystemMap {
		// First, load layouts into the common template
		layoutPath := LayoutsDir + "/*" + tm.extension
		_, err := commonTemplates.ParseFS(fsys, layoutPath)
		if err != nil {
			return nil, err
		}

		processPartials := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && filepath.Ext(path) == tm.extension {
				fullPath := path

				// Parse the partial template in the common template
				_, err := commonTemplates.ParseFS(fsys, fullPath)
				if err != nil {
					return err
				}

				//layoutPath := LayoutsDir + "/*" + tm.extension
				//_, err := commonTemplates.ParseFS(fsys, layoutPath, fullPath)
				//
				//if err != nil {
				//	return err
				//}
			}
			return nil
		}

		// If the "partials" directory exists, parse it
		if _, err := fsys.Open(PartialsDir); err == nil {
			if err := fs.WalkDir(fsys, PartialsDir, processPartials); err != nil {
				return nil, err
			}
		}
	}

	return commonTemplates, nil
}

func (tm *TemplateManager) printTemplateNames() {
	for name, tmpl := range tm.templates {
		tm.logger.Info("Template", slog.String("name", name))
		associatedTemplates := tmpl.Templates()
		for _, tmpl := range associatedTemplates {
			tm.logger.Info("    Partial/Child", slog.String("name", tmpl.Name()))
		}
	}
}

func (tm *TemplateManager) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tm *TemplateManager) render(w http.ResponseWriter, r *http.Request, resp *Response) {
	//path := tm.pathWithExtension(resp.TemplatePath())
	path := resp.TemplatePath()
	tmpl, ok := tm.templates[path]
	if !ok {
		tm.handleError(w, r, fmt.Errorf("%w: %s", ErrTempNotFound, resp.TemplatePath()))
		return
	}

	// Creating a buffer, so we can capture write errors before we write to the header
	// Note that layouts are always defined with the same name as the layout file without the extension (e.g. base.html -> base)
	buf := new(bytes.Buffer)
	layout := fmt.Sprintf("layout:%s", resp.TemplateLayout())
	err := tmpl.ExecuteTemplate(buf, layout, resp.ViewData(r).Data())
	if err != nil {
		path := tm.viewsPath(SystemDir, "server-error")
		if resp.TemplatePath() == path {
			http.Error(w, fmt.Errorf("error executing template: %w", err).Error(), http.StatusInternalServerError)
		} else {
			tm.handleError(w, r, fmt.Errorf("error executing template: %w", err))
		}
		return
	}

	// Add any additional headers
	for key, value := range resp.Headers() {
		w.Header().Set(key, value)
	}

	// Set the status code
	w.WriteHeader(resp.StatusCode())

	// Write the buffer to the response
	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tm *TemplateManager) viewsPath(path ...string) string {
	// For each path, append to the ViewsDir, separated by a slash
	return fmt.Sprintf("%s/%s", ViewsDir, strings.Join(path, "/"))
}

// pathWithExtension returns the path for the page template with the appropriate extension added.
func (tm *TemplateManager) pathWithExtension(path string) string {
	// Clean the path and add the extension
	curPath := strings.TrimSpace(path)

	// If the path is empty, return the default page path
	if curPath == "" {
		return fmt.Sprintf("home.%s", tm.extension)
	}

	// If the path does not have an extension, add the configured extension
	if filepath.Ext(curPath) == "" {
		return fmt.Sprintf("%s%s", curPath, tm.extension)
	}

	return curPath
}
