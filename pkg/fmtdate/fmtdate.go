package fmtdate

import (
	"strconv"
	"time"
)

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

// MySQLDateNow returns UTC date now in MySql format
func MySQLDateNow() string {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc).Format("2006-01-02")
}

// Timespan converts a MySql Datetime string to
// a timespan format ("how much time ago?").
func Timespan(mysqlDatetime string) string {
	ti, _ := time.Parse("2006-01-02 15:04:05", mysqlDatetime)
	secs := int(time.Since(ti).Seconds())
	mins := int(time.Since(ti).Minutes())
	hours := int(time.Since(ti).Hours())
	days := int(hours / 24)
	weeks := int(hours / 24 / 7)
	years := int(hours / 24 / 365)

	if years > 0 {
		suffix := "year"
		if years > 1 {
			suffix = "years"
		}
		return strconv.Itoa(years) + " " + suffix + " ago"
	} else if weeks > 0 {
		suffix := "week"
		if weeks > 1 {
			suffix = "weeks"
		}
		return strconv.Itoa(weeks) + " " + suffix + " ago"
	} else if days > 0 {
		suffix := "day"
		if days > 1 {
			suffix = "days"
		}
		return strconv.Itoa(days) + " " + suffix + " ago"
	} else if hours > 0 {
		suffix := "hour"
		if hours > 1 {
			suffix = "hours"
		}
		return strconv.Itoa(hours) + " " + suffix + " ago"
	} else if mins > 0 {
		suffix := "minute"
		if mins > 1 {
			suffix = "minutes"
		}
		return strconv.Itoa(mins) + " " + suffix + " ago"
	}
	suffix := "second"
	if secs > 1 {
		suffix = "seconds"
	}
	return strconv.Itoa(secs) + " " + suffix + " ago"
}
