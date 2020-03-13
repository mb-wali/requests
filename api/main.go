package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// API defines REST API of the requests service.
type API struct {
	Echo    *echo.Echo
	Title   string
	Version string
}

// RootResponse describes the response of the root endpoint.
type RootResponse struct {

	// The name of the service.
	Service string `json:"service"`

	// The service title.
	Title string `json:"title"`

	// The service version
	Version string `json:"Version"`
}

// ErrorResponse describes an error response for any endpoint.
type ErrorResponse struct {
	Description string `json:"description"`
}

// RootHandler handles GET requests to the / endpoint.
func (a API) RootHandler(ctx echo.Context) error {
	resp := RootResponse{
		Service: "requests",
		Title:   a.Title,
		Version: a.Version,
	}
	return ctx.JSON(http.StatusOK, resp)
}
