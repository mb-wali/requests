package api

import (
	"net/http"

	"github.com/cyverse-de/requests/db"
	"github.com/cyverse-de/requests/model"
	"github.com/labstack/echo"
)

// GetRequestTypesHandler handles GET requests to the /request-types endpoint.
func (a *API) GetRequestTypesHandler(ctx echo.Context) error {

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	// Obtain the list of request types.
	requestTypes, err := db.ListRequestTypes(tx)
	if err != nil {
		return err
	}

	// Build the response.
	resp := model.RequestTypeListing{RequestTypes: requestTypes}
	return ctx.JSON(http.StatusOK, resp)
}
