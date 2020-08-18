package fmtcsv

import (
	"errors"
	"strconv"
	"strings"
)

// MakeCSVFromList turns an int64 list into a csv string
func MakeCSVFromList(nums []int64) (string, error) {
	if len(nums) == 0 {
		return "", errors.New("Empty CSV")
	}
	var b strings.Builder
	for _, v := range nums {
		b.WriteString(strconv.FormatInt(v, 10) + ",")
	}
	csv := b.String()
	csv = csv[0 : len(csv)-1]
	return csv, nil
}
