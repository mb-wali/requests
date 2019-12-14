package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// API defines REST API of the requests service.
type API struct {
	Echo *echo.Echo
}

// RootResponse describes the response of the root endpoint.
type RootResponse struct {
	Service string
}

// ErrorResponse describes an error response for any endpoint.
type ErrorResponse struct {
	Description string
}

// RootHandler handles GET requests to the / endpoint.
func (a API) RootHandler(ctx echo.Context) error {
	resp := RootResponse{
		Service: "requests",
	}
	return ctx.JSON(http.StatusOK, resp)
}
