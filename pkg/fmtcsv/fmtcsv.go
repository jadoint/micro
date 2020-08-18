package fmtcsv

import (
	"strconv"
	"strings"
)

// MakeCSVFromList turns an int64 list into a csv string
func MakeCSVFromList(nums []int64) string {
	if len(nums) == 0 {
		return ""
	}
	var b strings.Builder
	for _, v := range nums {
		b.WriteString(strconv.FormatInt(v, 10) + ",")
	}
	csv := b.String()
	csv = csv[0 : len(csv)-1]
	return csv
}
