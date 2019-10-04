package fmtdate

import "time"

// MySQLUTC formats time to MySql datetime
func MySQLUTC(t time.Time) string {
	loc, _ := time.LoadLocation("UTC")
	return t.In(loc).Format("2006-01-02 15:04:05")
}

// MySQLUTCNow returns UTC date and time now in MySql format
func MySQLUTCNow() string {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}
