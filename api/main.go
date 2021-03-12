package api

import (
	"database/sql"
	"net/http"

	"github.com/cyverse-de/requests/clients/iplantemail"
	"github.com/cyverse-de/requests/clients/iplantgroups"
	"github.com/cyverse-de/requests/clients/notificationagent"
	"github.com/labstack/echo"
)

// API defines REST API of the requests service.
type API struct {
	Echo                    *echo.Echo
	Title                   string
	Version                 string
	DB                      *sql.DB
	UserDomain              string
	AdminEmail              string
	IPlantEmailClient       *iplantemail.Client
	IPlantGroupsClient      *iplantgroups.Client
	NotificationAgentClient *notificationagent.Client
}

// RootResponse describes the response of the root endpoint.
type RootResponse struct {

	// The name of the service
	Service string `json:"service"`

	// The service title
	Title string `json:"title"`

	// The service version
	Version string `json:"Version"`
}

// ErrorResponse describes an error response for any endpoint.
type ErrorResponse struct {
	Message string `json:"message"`
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
