package gtml

type gtmlError string

func (e gtmlError) Error() string {
	return string(e)
}

const (
	// ErrTempNotFound is returned when a template is not found.
	ErrTempNotFound = gtmlError("template not found")

	// ErrTempParse is returned when a template cannot be parsed.
	ErrTempParse = gtmlError("template parse error")

	// ErrTempRender is returned when a template cannot be rendered.
	ErrTempRender = gtmlError("template render error")
)
