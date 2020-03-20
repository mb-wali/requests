package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/cyverse-de/requests/db"
	"github.com/cyverse-de/requests/model"
	"github.com/cyverse-de/requests/query"

	"github.com/labstack/echo"
)

// updateRequestStatus updates the status of a request.
func (a *API) updateRequestStatus(tx *sql.Tx, requestID, status, userID, message string) error {
	var err error

	// Look up the request status code.
	requestStatusCode, err := db.GetRequestStatusCode(tx, status)
	if err != nil {
		return err
	}
	if requestStatusCode == nil {
		return fmt.Errorf("request status code not found: %s", status)
	}

	// Store the status request status update in the database.
	return db.AddRequestStatusUpdate(tx, requestID, requestStatusCode.ID, userID, message)
}

// AddRequestHandler handles POST requests to the /requests endpoint.
func (a *API) AddRequestHandler(ctx echo.Context) error {
	var err error

	// Extract and validate the user query parameter.
	user, err := query.ValidatedQueryParam(ctx, "user", "required")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "missing required query parameter: user",
		})
	}

	// Extract and validate the request body.
	requestSubmission := new(model.RequestSubmission)
	if err = ctx.Bind(requestSubmission); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("invalid request body: %s", err.Error()),
		})
	}
	if err = ctx.Validate(requestSubmission); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("invalid reuqest body: %s", err.Error()),
		})
	}

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	// Look up the user ID.
	userID, err := db.GetUserID(tx, user, a.UserDomain)
	if err != nil {
		tx.Rollback()
		return err
	}
	if userID == "" {
		tx.Rollback()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("user not found in DE database: %s", user),
		})
	}

	// Look up the request type.
	requestType, err := db.GetRequestType(tx, requestSubmission.RequestType)
	if err != nil {
		tx.Rollback()
		return err
	}
	if requestType == nil {
		tx.Rollback()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("request type not found: %s", requestSubmission.RequestType),
		})
	}

	// Store the request in the database.
	requestID, err := db.AddRequest(tx, userID, requestType.ID, requestSubmission.Details)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Store the request update in the database.
	err = a.updateRequestStatus(tx, requestID, "submitted", userID, "Request submitted.")
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

	// Extract the request body.
	return ctx.JSON(http.StatusOK, model.RequestSummary{
		ID:             requestID,
		RequestingUser: user,
		RequestType:    requestSubmission.RequestType,
		Details:        requestSubmission.Details,
	})
}
