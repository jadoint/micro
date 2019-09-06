package now

import "time"

// MySQLUTC returns UTC date and time now in MySql format
func MySQLUTC() string {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}
