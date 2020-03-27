package db

import (
	"database/sql"
	"fmt"

	"github.com/cyverse-de/requests/model"
)

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

// GetRequestStatusUpdate returns information for the request with the given ID.
func GetRequestStatusUpdate(tx *sql.Tx, updateID string) (*model.RequestUpdate, error) {
	query := `SELECT ru.id, rsc.name, regexp_replace(u.username, '@.*', ''), ru.created_date, ru.message
		FROM request_updates ru
		JOIN request_status_codes rsc ON ru.request_status_code_id = rsc.id
		JOIN users u ON ru.updating_user_id = u.id
		WHERE ru.id = $1`

	// Query the database.
	rows, err := tx.Query(query, updateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Quit now if no update was found.
	if !rows.Next() {
		return nil, nil
	}

	// Extract the status update information.
	var update model.RequestUpdate
	err = rows.Scan(&update.ID, &update.StatusCode, &update.UpdatingUser, &update.CreatedDate, &update.Message)
	if err != nil {
		return nil, err
	}
	return &update, nil
}

// AddRequestStatusUpdate adds a status update to a request.
func AddRequestStatusUpdate(
	tx *sql.Tx, requestID, requestStatusCodeID, updatingUserID, message string,
) (*model.RequestUpdate, error) {
	query := `INSERT INTO request_updates (request_id, request_status_code_id, updating_user_id, message)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id`

	//  Insert the request update.
	rows, err := tx.Query(query, requestID, requestStatusCodeID, updatingUserID, message)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// There should be at least one row in the result set.
	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned after an insert returning a value")
	}

	// Extract the request update id.
	var updateID string
	err = rows.Scan(&updateID)
	if err != nil {
		return nil, err
	}

	// The rows have to be closed before we can make additional queries.
	rows.Close()

	// Look up the update information.
	updateDetails, err := GetRequestStatusUpdate(tx, updateID)
	if err != nil {
		return nil, err
	}

	// The update should really exist since we just inserted it.
	if updateDetails == nil {
		return nil, fmt.Errorf("unable to look up the update that was just inserted")
	}

	// Return the update details.
	return updateDetails, nil
}
