package htmx

import (
	"net/http"
)

// IsHtmxRequest returns true if the current request is an HX/HTMX request, but not boosted. This distinguishes between
// a standard HTMX request and a boosted request as the server may need to differentiate between the two. HTMX sends
// a HX-Request header along with HX-Boosted, so it is not enough to check the HX-Request header.
func IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get(HXRequest) == "true" && r.Header.Get(HXBoosted) != "true"
}

// IsBoostedRequest returns true if the current request is boosted by the HX-Boosted header.
func IsBoostedRequest(r *http.Request) bool {
	return r.Header.Get(HXBoosted) == "true"
}

// IsAnyHtmxRequest returns true if the current request is either an HX/HTMX request or boosted.
func IsAnyHtmxRequest(r *http.Request) bool {
	return IsHtmxRequest(r) || IsBoostedRequest(r)
}

// IsHistoryRestoreRequest returns true if the current request contains the HX-History-Restore header
func IsHistoryRestoreRequest(r *http.Request) bool {
	return r.Header.Get(HXHistoryRestoreRequest) == "true"
}

// CurrentURL returns the HX-Current-URL header, if it exists
func CurrentURL(r *http.Request) (string, bool) {
	if _, ok := r.Header[http.CanonicalHeaderKey(HXCurrentURL)]; !ok {
		return "", false
	}

	return r.Header.Get(HXCurrentURL), true
}

// Prompt returns the HX-Prompt header, if it exists
func Prompt(r *http.Request) (string, bool) {
	if _, ok := r.Header[http.CanonicalHeaderKey(HXPrompt)]; !ok {
		return "", false
	}

	return r.Header.Get(HXPrompt), true
}

// Target returns the HX-Target header, if it exists
func Target(r *http.Request) (string, bool) {
	if _, ok := r.Header[http.CanonicalHeaderKey(HXTarget)]; !ok {
		return "", false
	}

	return r.Header.Get(HXTarget), true
}

// Trigger returns the HX-Trigger header, if it exists
func Trigger(r *http.Request) (string, bool) {
	if _, ok := r.Header[http.CanonicalHeaderKey(HXTrigger)]; !ok {
		return "", false
	}

	return r.Header.Get(HXTrigger), true
}

// TriggerName returns the HX-Trigger-Name header, if it exists
func TriggerName(r *http.Request) (string, bool) {
	if _, ok := r.Header[http.CanonicalHeaderKey(HXTriggerName)]; !ok {
		return "", false
	}

	return r.Header.Get(HXTriggerName), true
}
