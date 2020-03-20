package db

import (
	"database/sql"
	"fmt"
)

// GetUserID obtains the internal user ID  for the given username and user domain from the DE database.
func GetUserID(tx *sql.Tx, username, userDomain string) (string, error) {
	query := "SELECT id FROM users WHERE username = $1"
	qualifiedUsername := fmt.Sprintf("%s@%s", username, userDomain)

	// Query the database.
	rows, err := tx.Query(query, qualifiedUsername)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// extract the user ID if at least one row was found.
	var userID string
	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return "", err
		}
	}

	return userID, nil
}
