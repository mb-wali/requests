package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cyverse-de/requests/model"

	sq "github.com/Masterminds/squirrel"

	"github.com/pkg/errors"
)

// AddRequest adds a new request to the database.
func AddRequest(tx *sql.Tx, userID, requestTypeID string, details interface{}) (string, error) {
	query := `INSERT INTO requests (request_type_id, requesting_user_id, details)
			  VALUES ($1, $2, CAST($3 AS json))
			  RETURNING id`

	// Encode the request details.
	encodedDetails, err := json.Marshal(details)
	if err != nil {
		return "", errors.Wrap(err, "unable to JSON encode the request details")
	}

	// Insert the new request.
	rows, err := tx.Query(query, requestTypeID, userID, encodedDetails)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// There should be at least one row in the result set.
	if !rows.Next() {
		return "", fmt.Errorf("no rows returned after an insert returning a value")
	}

	// Extract the request ID.
	var requestID string
	err = rows.Scan(&requestID)
	if err != nil {
		return "", err
	}

	return requestID, nil
}

// RequestListingOptions represents options that can be used to filter a request listing.
type RequestListingOptions struct {
	IncludeCompletedRequests bool
	RequestType              string
}

// GetRequestListing obtains a list of requests from the database.
func GetRequestListing(tx *sql.Tx, options *RequestListingOptions) ([]*model.RequestSummary, error) {
	base := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("r.id, regexp_replace(u.username, '@.*', ''), rt.name, r.details").
		From("requests r").
		Join("users u ON r.requesting_user_id = u.id").
		Join("request_types rt ON r.request_type_id = rt.id")

	// Add the filter to omit completed requests if we're not supposed to include them in the listing.
	if !options.IncludeCompletedRequests {
		nestedBuilder := sq.StatementBuilder.
			Select("*").
			Prefix("NOT EXISTS (").
			From("request_updates ru").
			Join("request_status_codes rsc ON ru.request_status_code_id = rsc.id").
			Where("ru.request_id = r.id").
			Where(sq.Eq{"rsc.name": []string{"complete", "rejected"}}).
			Suffix(")")
		base = base.Where(nestedBuilder)
	}

	// Add the filter to limit the listing to requests of a given type if applicable.
	if options.RequestType != "" {
		base = base.Where(sq.Eq{"rt.name": options.RequestType})
	}

	// Build the query.
	query, args, err := base.ToSql()
	if err != nil {
		return nil, err
	}

	// Query the database.
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build the listing from the result set.
	listing := make([]*model.RequestSummary, 0)
	for rows.Next() {
		var request model.RequestSummary
		var requestDetails string
		err = rows.Scan(&request.ID, &request.RequestingUser, &request.RequestType, &requestDetails)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(requestDetails), &request.Details)
		if err != nil {
			return nil, err
		}
		listing = append(listing, &request)
	}

	return listing, nil
}

// GetRequestStatusUpdates looks up the status updates for a request.
func GetRequestStatusUpdates(tx *sql.Tx, requestID string) ([]*model.RequestUpdate, error) {
	query := `SELECT ru.id, rsc.name, regexp_replace(u.username, '@.*', ''), ru.created_date, ru.message
			  FROM request_updates ru
			  JOIN request_status_codes rsc ON ru.request_status_code_id = rsc.id
			  JOIN users u ON ru.updating_user_id = u.id
			  WHERE ru.request_id = $1
			  ORDER BY ru.created_date`

	// Query the database.
	rows, err := tx.Query(query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build the array of status updates.
	updates := make([]*model.RequestUpdate, 0)
	for rows.Next() {
		var update model.RequestUpdate
		err := rows.Scan(&update.ID, &update.StatusCode, &update.UpdatingUser, &update.CreatedDate, &update.Message)
		if err != nil {
			return nil, err
		}
		updates = append(updates, &update)
	}
	return updates, nil
}

// GetRequestDetails looks up the details of a request.
func GetRequestDetails(tx *sql.Tx, id string) (*model.RequestDetails, error) {
	query := `SELECT r.id, regexp_replace(u.username, '@.*', ''), rt.name, r.details
			  FROM requests r
			  JOIN users u ON r.requesting_user_id = u.id
			  JOIN request_types rt ON r.request_type_id = rt.id
			  WHERE r.id = $1`

	// Query the database.
	rows, err := tx.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Just return nil if there aren't any rows.
	if !rows.Next() {
		return nil, nil
	}

	// Extract the request details.
	var rd model.RequestDetails
	var rdDetails string
	err = rows.Scan(&rd.ID, &rd.RequestingUser, &rd.RequestType, &rdDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(rdDetails), &rd.Details)
	if err != nil {
		return nil, err
	}

	// The rows have to be closed before we can make additional queries.
	rows.Close()

	// Add status information to the request details.
	rd.Updates, err = GetRequestStatusUpdates(tx, id)
	if err != nil {
		return nil, err
	}

	return &rd, nil
}
