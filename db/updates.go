package db

import "database/sql"

// AddRequestStatusUpdate adds a request status update record to the database.
func AddRequestStatusUpdate(tx *sql.Tx, requestID, requestStatusCodeID, updatingUserID, message string) error {
	stmt := `INSERT INTO request_updates (request_id, request_status_code_id, updating_user_id, message)
			 VALUES ($1, $2, $3, $4)`

	// Execute the statement.
	_, err := tx.Exec(stmt, requestID, requestStatusCodeID, updatingUserID, message)

	return err
}
