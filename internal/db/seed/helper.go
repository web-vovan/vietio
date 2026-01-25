package seed

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// список всех uuid в таблице
func getAllUuidFromTable(dbConn *sql.DB, table string) ([]uuid.UUID, error) {
	query := fmt.Sprintf("SELECT uuid FROM %s", table)

	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []uuid.UUID

	for rows.Next() {
		var uuid uuid.UUID
		if err := rows.Scan(&uuid); err != nil {
			return nil, err
		}

		result = append(result, uuid)
	}

	return result, nil
}

// список всех id в таблице
func getAllIdsFromTable(dbConn *sql.DB, table string) ([]int64, error) {
	query := fmt.Sprintf("SELECT id FROM %s", table)

	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []int64

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		result = append(result, id)
	}

	return result, nil
}
