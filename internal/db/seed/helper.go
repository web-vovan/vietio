package seed

import (
	"database/sql"
	"fmt"
)

// список всех id в таблице
func getAllIdsFromTable(dbConn *sql.DB, table string) ([]int64, error) {
    query := fmt.Sprintf("SELECT id FROM %s", table)
    
    rows, err := dbConn.Query(query)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var ids []int64

    for rows.Next() {
        var id int64
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }

        ids = append(ids, id)
    }

    return ids, nil
}