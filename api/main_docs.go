// Package api DE Administrative Requests API
//
// Documentation of the DE Administrative Requests API
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package api

// swagger:route GET / misc getRoot
// Returns general information about the API.
// responses:
//   200: rootResponse

// General information about hte API.
// swagger:response rootResponse
type rootResponseWrapper struct {
	// in:body
	Body RootResponse
}

// Basic error response.
// swagger:response errorResponse
type errorResponseWrapper struct {
	Body ErrorResponse
}
