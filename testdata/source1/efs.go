package source1

import "embed"

//go:embed "layouts" "partials" "views"
var FS embed.FS
