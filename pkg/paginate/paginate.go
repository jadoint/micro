package paginate

import (
	"net/http"
	"strconv"
)

// GetPageNum gets pageNum from URL query parameter
func GetPageNum(r *http.Request) (int, error) {
	pageNumParam := r.URL.Query().Get("pageNum")
	var err error
	pageNum := 1
	if pageNumParam != "" {
		pageNum, err = strconv.Atoi(pageNumParam)
		if err != nil {
			return 0, err
		}
		if pageNum <= 0 {
			pageNum = 1
		}
	}
	return pageNum, nil
}

// GetOffset translates pageNum from a URL query
// parameter into an offset usable by a LIMIT clause.
func GetOffset(pageNum int, pageSize int) int {
	return (pageNum - 1) * pageSize
}
