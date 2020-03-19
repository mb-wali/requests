package db

import (
	"database/sql"

	"github.com/cyverse-de/requests/model"
)

// requestStatusCodesFromRows converts SQL rows to an array of request status codes. The column order is always
// expected to be ID, name, display_name, and email_template.
func requestStatusCodesFromRows(rows *sql.Rows) ([]*model.RequestStatusCode, error) {
	requestStatusCodes := make([]*model.RequestStatusCode, 0)

	// Build the list of request status codes.
	for rows.Next() {
		var rsc model.RequestStatusCode
		err := rows.Scan(&rsc.ID, &rsc.Name, &rsc.DisplayName, &rsc.EmailTemplate)
		if err != nil {
			return nil, err
		}
		requestStatusCodes = append(requestStatusCodes, &rsc)
	}

	return requestStatusCodes, nil
}

// ListRequestStatusCodes lists all of the currently available request status codes.
func ListRequestStatusCodes(tx *sql.Tx) ([]*model.RequestStatusCode, error) {
	query := "SELECT id, name, display_name, email_template FROM request_status_codes"

	// Query the database.
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Return the list of rows.
	return requestStatusCodesFromRows(rows)
}
