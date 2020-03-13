package db

import (
	"database/sql"

	"github.com/cyverse-de/requests/model"
)

func requestTypesFromRows(rows *sql.Rows) ([]*model.RequestType, error) {
	requestTypes := make([]*model.RequestType, 0)

	// Build the list of request types.
	for rows.Next() {
		var rt model.RequestType
		err := rows.Scan(&rt.ID, &rt.Name)
		if err != nil {
			return nil, err
		}
		requestTypes = append(requestTypes, &rt)
	}

	return requestTypes, nil
}

// ListRequestTypes returns a listing of request types from the database sorted by name.
func ListRequestTypes(tx *sql.Tx) ([]*model.RequestType, error) {
	query := "SELECT id, name FROM request_types ORDER BY name"

	// Query the database.
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Return the list of request types.
	return requestTypesFromRows(rows)
}
