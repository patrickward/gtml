package gtml

import (
	"fmt"
	"net/http"

	"github.com/patrickward/gtml/htmx"
	"github.com/patrickward/gtml/htmx/location"
	"github.com/patrickward/gtml/htmx/swap"
)

// HxLocation sets the HX-Location header, which instructs the browser to navigate to the given path without reloading the page.
//
// For simple navigations, use a path and no HxLocation options.
// For more complex navigations, use the HxLocation options to fine-tune the navigation.
//
// Sets the HX-Location header with the given path.
//
// For more information, see: https://htmx.org/headers/hx-location
func (resp *Response) HxLocation(path string, opt ...location.Option) *Response {
	if opt == nil {
		resp.headers[htmx.HXLocation] = path
	} else {
		loc := location.NewLocation(path, opt...)
		resp.headers[htmx.HXLocation] = loc.String()
	}
	return resp
}

// HxPushURL sets the HX-Push-Url header, which instructs the browser to navigate to the given path without reloading the page.
//
// To prevent the browser from updating the page, set the HX-Push-Url header to an empty string or "false".
//
// For more information, see: https://htmx.org/headers/hx-push-url
func (resp *Response) HxPushURL(path string) *Response {
	resp.headers[htmx.HXPushURL] = path
	return resp
}

// HxNoPushURL prevents the browser from updating the history stack by setting the HX-Push-Url header to "false".
//
// For more information, see: https://htmx.org/headers/hx-no-push-url
func (resp *Response) HxNoPushURL() *Response {
	resp.headers[htmx.HXPushURL] = "false"
	return resp
}

// HxRedirect sets the HX-Redirect header, which instructs the browser to navigate to the given path (this will reload the page).
//
// For more information, see: https://htmx.org/reference/#response_headers
func (resp *Response) HxRedirect(path string) *Response {
	resp.headers[htmx.HXRedirect] = path
	return resp
}

// HxRefresh sets the HX-Refresh header, which instructs the browser to reload the page.
//
// For more information, see: https://htmx.org/reference/#response_headers
func (resp *Response) HxRefresh() *Response {
	resp.headers[htmx.HXRefresh] = "true"
	return resp
}

// HxNoRefresh prevents the browser from reloading the page by setting the HX-Refresh header to "false".
//
// For more information, see: https://htmx.org/reference/#response_headers
func (resp *Response) HxNoRefresh() *Response {
	resp.headers[htmx.HXRefresh] = "false"
	return resp
}

// HxReplaceURL sets the HX-Replace-Url header, which instructs the browser to replace the history stack with the given path.
//
// For more information, see: https://htmx.org/headers/hx-replace-url
func (resp *Response) HxReplaceURL(path string) *Response {
	resp.headers[htmx.HXReplaceURL] = path
	return resp
}

// HxNoReplaceURL prevents the browser from updating the history stack by setting the HX-Replace-Url header to "false".
//
// For more information, see: https://htmx.org/headers/hx-replace-url
func (resp *Response) HxNoReplaceURL() *Response {
	resp.headers[htmx.HXReplaceURL] = "false"
	return resp
}

// HxReswap sets the HX-Reswap header, which instructs HTMX to change the swap behavior of the target element.
//
// For more information, see: https://htmx.org/attributes/hx-swap
func (resp *Response) HxReswap(swap *swap.Style) *Response {
	resp.headers[htmx.HXReswap] = swap.String()
	return resp
}

// HxRetarget sets the HX-Retarget header, which instructs HTMX to update the target element.
//
// For more information, see: https://htmx.org/reference/#response_headers
func (resp *Response) HxRetarget(target string) *Response {
	resp.headers[htmx.HXRetarget] = target
	return resp
}

// HxReselect sets the HX-ReSelect header, which instructs HTMX to update which part of the response is selected.
//
// For more information, see: https://htmx.org/reference/#response_headers and https://htmx.org/attributes/hx-select
func (resp *Response) HxReselect(reselect string) *Response {
	resp.headers[htmx.HXReselect] = reselect
	return resp
}

// HxTrigger sets a HX-Trigger header
//
// For more information, see: https://htmx.org/headers/hx-trigger/
func (resp *Response) HxTrigger(event string, value any) *Response {
	resp.triggers.Set(event, value)
	return resp
}

// HxTriggerAfterSettle sets a HX-Trigger-After-Settle header
//
// For more information, see: https://htmx.org/headers/hx-trigger/
func (resp *Response) HxTriggerAfterSettle(event string, value any) *Response {
	resp.triggers.SetAfterSettle(event, value)
	return resp
}

// HxTriggerAfterSwap sets a HX-Trigger-After-Swap header
//
// For more information, see: https://htmx.org/headers/hx-trigger/
func (resp *Response) HxTriggerAfterSwap(event string, value any) *Response {
	resp.triggers.SetAfterSwap(event, value)
	return resp
}

// HxNonce returns the HTMX nonce value from the request context, if available.
// This adds the inlineScriptNonce key to a JSON object with the nonce value and can be used in an HTMX meta tag.
func (v *ResponseData) HxNonce() string {
	return fmt.Sprintf("{\"includeIndicatorStyles\":false,\"inlineScriptNonce\": \"%s\"}", v.Nonce())
}

// HxLayout sets the layout for HTMX requests if the request is an HTMX request, otherwise it uses the default layout.
//
// Parameters:
//   - request is used to determine if the request is an HTMX request.
//   - hxLayout is the layout to use for HTMX requests.
//   - layout is the default layout to use if the request is not an HTMX request.
func (resp *Response) HxLayout(r *http.Request, hxLayout, layout string) *Response {
	if htmx.IsHtmxRequest(r) {
		resp.Layout(hxLayout)
	} else {
		resp.Layout(layout)
	}
	return resp
}

// IsHtmxRequest returns true if the request is an HTMX request, but not a boosted request.
func (v *ResponseData) IsHtmxRequest() bool {
	return htmx.IsHtmxRequest(v.request)
}

// IsBoostedRequest returns true if the request is a boosted request.
func (v *ResponseData) IsBoostedRequest() bool {
	return htmx.IsBoostedRequest(v.request)
}
