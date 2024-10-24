# gtml - Experimental

Template Manager, response service, and helpers for my Go web projects

This is a simple template manager for Go web projects. It provides a way to manage templates, render them, and serve them as HTTP responses.

### Basic Usage

To use the template manager, create a new instance of `TemplateManager` with the desired options. Then, you can use the manager to create responses and render templates.

```go
// Create a new template manager
tm, err := gtml.NewTemplateManager(
	gtml.Source{"-": &templates.Files},
    gtml.TemplateManagerOptions{
        Extension: ".gtml",
        Funcs:  funcs.TemplateFuncs,
        Logger: logger,
    })

if err != nil {
    return fmt.Errorf("error initializing template manager: %w", err)
}

// Then, you might want to set it on a server or handler
server := NewServer(tm)

// Later, create a new response and render a template
data := h.server.NewTemplateData(r)

h.server.TM().NewResponse().
    Title("Home").
    Path("home").
    Data(data).
    Render(w, r)
```

The `TemplateManager` handles loading templates from the specified sources, rendering them with the provided data, and serving them as HTTP responses. It will also handle any errors that occur during the rendering process.

