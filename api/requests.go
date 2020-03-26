package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyverse-de/requests/db"
	"github.com/cyverse-de/requests/model"
	"github.com/cyverse-de/requests/query"

	"github.com/labstack/echo"
)

// copyRequestDetails makes a one-level-deep copy of a map. For copying request details, we only need to go one level
// deep because this service doesn't need to modify anything below the top level of the map.
func copyRequestDetails(requestDetails map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for k, v := range requestDetails {
		copy[k] = v
	}
	return copy
}

// formatRequestDetails builds a human readable representation of a set of request details. For now, we're just going
// to turn it into a pretty-printed JSON document.
func formatRequestDetails(requestDetails map[string]interface{}) (string, error) {
	formattedDetails, err := json.MarshalIndent(requestDetails, "", "  ")
	if err != nil {
		return "", err
	}
	return string(formattedDetails), nil
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

	// Look up the request status code.
	requestStatusCode, err := db.GetRequestStatusCode(tx, "submitted")
	if err != nil {
		tx.Rollback()
		return err
	}
	if requestStatusCode == nil {
		tx.Rollback()
		return fmt.Errorf("request status code not found: submitted")
	}

	// Store the request update in the database.
	err = db.AddRequestStatusUpdate(tx, requestID, requestStatusCode.ID, userID, "Request submitted.")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Format a human readable copy of the request submission details.
	humanReadableRequestDetails, err := formatRequestDetails(requestSubmission.Details.(map[string]interface{}))
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add required information to a copy of the request details.
	requestDetails := copyRequestDetails(requestSubmission.Details.(map[string]interface{}))
	requestDetails["username"] = user
	requestDetails["request_type"] = requestType.Name
	requestDetails["request_details"] = humanReadableRequestDetails

	// Send the email.
	err = a.IPlantEmailClient.SendRequestSubmittedEmail(a.AdminEmail, requestStatusCode.EmailTemplate, requestDetails)
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

	// Build the response body.
	return ctx.JSON(http.StatusOK, model.RequestSummary{
		ID:             requestID,
		RequestingUser: user,
		RequestType:    requestSubmission.RequestType,
		Details:        requestSubmission.Details,
	})
}

// GetRequestsHandler handles GET requests to the /requests endpoint.
func (a *API) GetRequestsHandler(ctx echo.Context) error {
	var err error

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	// Extract and validate the user query parameter.
	defaultIncludeCompleted := false
	includeCompleted, err := query.ValidateBooleanQueryParam(ctx, "include-completed", &defaultIncludeCompleted)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: err.Error(),
		})
	}

	// Build the request listing obtions.
	options := &db.RequestListingOptions{
		IncludeCompletedRequests: includeCompleted,
	}

	// Get the list of matching requests.
	requests, err := db.GetRequestListing(tx, options)
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

	// Return the listing.
	return ctx.JSON(http.StatusOK, &model.RequestListing{
		Requests: requests,
	})
}

// GetRequestDetailsHandler handles GET requests to the /requests/:id endpoint.
func (a *API) GetRequestDetailsHandler(ctx echo.Context) error {
	id := ctx.Param("id")
	var err error

	// Start a transaction
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	// Look up the request details.
	requestDetails, err := db.GetRequestDetails(tx, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if requestDetails == nil {
		tx.Rollback()
		return ctx.JSON(http.StatusNotFound, &ErrorResponse{
			Message: fmt.Sprintf("request %s not found", id),
		})
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Return the response.
	return ctx.JSON(http.StatusOK, requestDetails)
}
