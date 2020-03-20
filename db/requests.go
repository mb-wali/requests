package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

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
