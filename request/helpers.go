package request

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// SchemeHostPort returns the scheme, host, and port for the given request. If behind a proxy and
// the X-Forwarded-Proto, X-Forwarded-Host, and/or X-Forwarded-Port headers are set, they will be used.
func SchemeHostPort(r *http.Request) (string, string, string) {
	var scheme, host, port string

	if r.Header.Get("X-Forwarded-Proto") != "" {
		scheme = r.Header.Get("X-Forwarded-Proto")
	} else if r.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}

	if r.Header.Get("X-Forwarded-Host") != "" {
		host = strings.Split(r.Header.Get("X-Forwarded-Host"), ":")[0]
	} else {
		host = strings.Split(r.Host, ":")[0]
	}

	if r.Header.Get("X-Forwarded-Port") != "" {
		port = r.Header.Get("X-Forwarded-Port")
	} else if strings.Contains(r.Host, ":") {
		_, port, _ = net.SplitHostPort(r.Host)
	} else if scheme == "https" {
		port = "443"
	} else {
		port = "80"
	}

	return scheme, host, port
}

// BaseURL returns the base URL for the given request.
func BaseURL(r *http.Request) string {
	scheme, host, port := SchemeHostPort(r)
	if scheme == "http" && port == "80" || scheme == "https" && port == "443" {
		return fmt.Sprintf("%s://%s", scheme, host)
	}
	return fmt.Sprintf("%s://%s:%s", scheme, host, port)
}

// Scheme returns the scheme (http or https) of the given request.
func Scheme(r *http.Request) string {
	scheme, _, _ := SchemeHostPort(r)
	return scheme
}

// Host returns the host part of the URL from the request.
func Host(r *http.Request) string {
	_, host, _ := SchemeHostPort(r)
	return host
}

// Port returns the port number from the HTTP request.
func Port(r *http.Request) string {
	_, _, port := SchemeHostPort(r)
	return port
}

// Method returns the HTTP method of the provided request.
func Method(r *http.Request) string {
	return r.Method
}

// URLPath returns the path portion of the URL in the given http.Request object.
func URLPath(r *http.Request) string {
	return r.URL.Path
}

func Referer(r *http.Request) string {
	return r.Referer()
}

// RemoteAddr returns the client's IP address extracted from the given http.Request. If the server is behind a proxy
// and the x-forwarded-for header is set, it will return the first IP address in the list. Otherwise, it will check
// for the x-real-ip header and return its value. If neither header is set, it will return the remote address of the request.
func RemoteAddr(r *http.Request) string {
	remoteAddr := ""
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ip := strings.Split(xForwardedFor, ",")[0]
		remoteAddr = strings.TrimSpace(ip)
	}

	if remoteAddr == "" {
		remoteAddr = r.Header.Get("X-Real-IP")
	}

	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	return remoteAddr
}

// IsJSONRequest returns true if the request has a Content-Type of application/json
func IsJSONRequest(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}

// IsFormRequest returns true if the request has a Content-Type of application/x-www-form-urlencoded
func IsFormRequest(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/x-www-form-urlencoded"
}

// IsXMLHttpRequest returns true if the request is an XMLHttpRequest (e.g. via fetch)
func IsXMLHttpRequest(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// IsSecure returns true if the current request is using the HTTPS scheme.
func IsSecure(r *http.Request) bool {
	scheme, _, _ := SchemeHostPort(r)
	return scheme == "https"
}

// UserAgent returns the value of the User-Agent header from the HTTP request.
func UserAgent(r *http.Request) string {
	return r.UserAgent()
}

// InPath returns true if the current request path matches the given path.
//
// Parameters:
// - r: *http.Request - the current HTTP request
// - path: string - the path to match against the request path
// - options: ...string - optional parameters for matching the path
//
// Returns:
// - bool: true if the request path matches the given path, false otherwise
//
// Options:
// - exact: matches the exact path
// - contains: matches if the path contains the given path
// - prefix: matches if the path starts with the given path
// - suffix: matches if the path ends with the given path
func InPath(r *http.Request, path string, options ...string) bool {
	requestPath := r.URL.Path
	if len(options) > 0 {
		option := strings.ToLower(options[0])
		if option == "exact" {
			if requestPath == "/" {
				return path == "/" || path == ""
			}
			return strings.Compare(requestPath, path) == 0
		} else if option == "contains" {
			return strings.Contains(requestPath, path)
		} else if option == "suffix" {
			return strings.HasSuffix(requestPath, path)
		} else {
			return strings.HasPrefix(requestPath, path)
		}
	} else {
		return strings.HasPrefix(requestPath, path)
	}
}
