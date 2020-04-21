package query

import (
	"database/sql"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/logger"
)

// BoolQueryHandler handles queries whose sole purpose is to
// determine if a query results in a row or not.
func BoolQueryHandler(clients *conn.Clients, sqlQuery string, args ...interface{}) bool {
	db := clients.DB.Read
	var isFound int
	err := db.QueryRow(sqlQuery, args...).Scan(&isFound)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(err)
		}
		return false
	}
	return true
}
