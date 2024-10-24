package htmx

// HTMX Request Headers
const (
	// HXBoosted indicates the request is an HTMX boosted request.
	HXBoosted = "HX-Boosted"

	// HXCurrentURL is the current URL of the browser
	HXCurrentURL = "HX-Current-URL"

	// HXHistoryRestoreRequest indicates the request is a history restore request due to a miss in the local history cache
	HXHistoryRestoreRequest = "HX-History-Restore-Request"

	// HXPrompt is the user response to an hx-prompt
	HXPrompt = "HX-Prompt"

	// HXRequest indicates the request is an HTMX request
	HXRequest = "HX-Request"

	// HXTarget is the ID of the target element, if any
	HXTarget = "HX-Target"

	// HXTriggerName is the name of the triggered element
	HXTriggerName = "HX-Trigger-Name"
)

// HTMX Response Headers
const (
	// HXLocation allows you to do a client-side redirect without a full page reload (similar to a boost)
	HXLocation = "HX-Location"

	// HXPushURL pushes the URL into the browser history stack
	HXPushURL = "HX-Push-URL"

	// HXRedirect allows you to do a client-side redirect as a full page reload
	HXRedirect = "HX-Redirect"

	// HXRefresh allows you to refresh the page as a full page reload
	HXRefresh = "HX-Refresh"

	// HXReplaceURL replaces the URL in the browser location bar
	HXReplaceURL = "HX-Replace-URL"

	// HXReswap allows you to specify how the response will be swapped in by HTMX
	HXReswap = "HX-Reswap"

	// HXRetarget is a CSS selector that updates the target of the content to update by HTMX
	HXRetarget = "HX-Retarget"

	// HXReselect is a CSS selector that indicates which part of the response should be swapped in. Overrides existing hx-swap values on the triggering element.
	HXReselect = "HX-ReSelect"

	// HXTriggerAfterSettle allows you to trigger client-side-events after the settle phase
	HXTriggerAfterSettle = "HX-Trigger-After-Settle"

	// HXTriggerAfterSwap allows you to trigger client-side-events after the swap phase
	HXTriggerAfterSwap = "HX-Trigger-After-Swap"
)

// Dual-use Headers (request and response)
const (
	// HXTrigger is the ID of the triggered element when used in a request. As a response header, it can be used to trigger client-side events.
	HXTrigger = "HX-Trigger"
)
