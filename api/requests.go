package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyverse-de/requests/clients/notificationagent"

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
	defer tx.Rollback()

	// Look up the user ID.
	userID, err := db.GetUserID(tx, user, a.UserDomain)
	if err != nil {
		return err
	}
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("user not found in DE database: %s", user),
		})
	}

	// Look up the request type.
	requestType, err := db.GetRequestType(tx, requestSubmission.RequestType)
	if err != nil {
		return err
	}
	if requestType == nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("request type not found: %s", requestSubmission.RequestType),
		})
	}

	// Store the request in the database.
	requestID, err := db.AddRequest(tx, userID, requestType.ID, requestSubmission.Details)
	if err != nil {
		return err
	}

	// Look up the request status code.
	requestStatusCode, err := db.GetRequestStatusCode(tx, "submitted")
	if err != nil {
		return err
	}
	if requestStatusCode == nil {
		return fmt.Errorf("request status code not found: submitted")
	}

	// Store the request update in the database.
	update, err := db.AddRequestStatusUpdate(tx, requestID, requestStatusCode.ID, userID, "Request submitted.")
	if err != nil {
		return err
	}

	// Format a human readable copy of the request submission details.
	humanReadableRequestDetails, err := formatRequestDetails(requestSubmission.Details.(map[string]interface{}))
	if err != nil {
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
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Build the response body.
	return ctx.JSON(http.StatusOK, model.RequestSummary{
		ID:             requestID,
		RequestingUser: user,
		RequestType:    requestSubmission.RequestType,
		CreatedDate:    update.CreatedDate,
		Status:         requestStatusCode.DisplayName,
		UpdatedDate:    update.CreatedDate,
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
	defer tx.Rollback()

	// Extract and validate the include-completed query parameter.
	defaultIncludeCompleted := false
	includeCompleted, err := query.ValidateBooleanQueryParam(ctx, "include-completed", &defaultIncludeCompleted)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	// Build the request listing obtions.
	options := &db.RequestListingOptions{
		IncludeCompletedRequests: includeCompleted,
		RequestType:              ctx.QueryParam("request-type"),
		RequestingUser:           ctx.QueryParam("requesting-user"),
	}

	// Get the list of matching requests.
	requests, err := db.GetRequestListing(tx, options)
	if err != nil {
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
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
	defer tx.Rollback()

	// Look up the request details.
	requestDetails, err := db.GetRequestDetails(tx, id)
	if err != nil {
		return err
	}
	if requestDetails == nil {
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("request %s not found", id),
		})
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Return the response.
	return ctx.JSON(http.StatusOK, requestDetails)
}

// UpdateRequestHandler handles POST requests to the /requests/:id/update endpoint.
func (a *API) UpdateRequestHandler(ctx echo.Context) error {
	id := ctx.Param("id")
	var err error

	// Extract and validate the user query parameter.
	user, err := query.ValidatedQueryParam(ctx, "user", "required")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "missing required query parameter: user",
		})
	}

	// Extract and validate the request body.
	requestUpdateSubmission := new(model.RequestUpdateSubmission)
	if err = ctx.Bind(requestUpdateSubmission); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("invalid request body: %s", err.Error()),
		})
	}
	if err = ctx.Validate(requestUpdateSubmission); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("invalid reuqest body: %s", err.Error()),
		})
	}

	// Start a transaction
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Look up the updating user ID.
	userID, err := db.GetUserID(tx, user, a.UserDomain)
	if err != nil {
		return err
	}
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("user not found in DE database: %s", user),
		})
	}

	// Verify that the request exists.
	request, err := db.GetRequestDetails(tx, id)
	if err != nil {
		return err
	}
	if request == nil {
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("request %s not found", id),
		})
	}

	// Look up the request status code.
	requestStatusCode, err := db.GetRequestStatusCode(tx, requestUpdateSubmission.StatusCode)
	if err != nil {
		return err
	}
	if requestStatusCode == nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("invalid request status code: %s", requestUpdateSubmission.StatusCode),
		})
	}

	// Save the request status update.
	update, err := db.AddRequestStatusUpdate(tx, id, requestStatusCode.ID, userID, requestUpdateSubmission.Message)
	if err != nil {
		return err
	}

	// Look up information about the user who submitted the request.
	requestingUserInfo, err := a.IPlantGroupsClient.GetUserInfo(request.RequestingUser)
	if err != nil {
		return err
	}

	// Format a human readable copy of the request submission details.
	humanReadableRequestDetails, err := formatRequestDetails(request.Details.(map[string]interface{}))
	if err != nil {
		return err
	}

	// Add required information to a copy of the request details.
	requestDetails := copyRequestDetails(request.Details.(map[string]interface{}))
	requestDetails["request_details"] = humanReadableRequestDetails
	requestDetails["update_message"] = update.Message
	requestDetails["email_address"] = requestingUserInfo.Email
	requestDetails["action"] = "request_status_change"
	requestDetails["user"] = requestingUserInfo.ID

	// Send the email.
	emailText := "Your administrative request status is now: " +
		requestStatusCode.DisplayName +
		"."
	err = a.NotificationAgentClient.SendNotification(
		&notificationagent.NotificationRequest{
			Type:          "request",
			User:          *requestingUserInfo.ID,
			Subject:       "Administrative Request " + requestStatusCode.DisplayName,
			Message:       emailText,
			Email:         true,
			EmailTemplate: requestStatusCode.EmailTemplate,
			Payload:       requestDetails,
		},
	)
	if err != nil {
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Return the response.
	return ctx.JSON(http.StatusOK, update)
}
