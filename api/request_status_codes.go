package api

import (
	"net/http"

	"github.com/cyverse-de/requests/db"
	"github.com/cyverse-de/requests/model"
	"github.com/labstack/echo"
)

// GetRequestStatusCodesHandler handles GET requests to the /request-status-codes endpoint.
func (a *API) GetRequestStatusCodesHandler(ctx echo.Context) error {

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	// Obtain the list of request status codes.
	requestStatusCodes, err := db.ListRequestStatusCodes(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Return the response.
	return ctx.JSON(http.StatusOK, model.RequestStatusCodeListing{
		RequestStatusCodes: requestStatusCodes,
	})
}
