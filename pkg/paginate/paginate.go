package paginate

import (
	"net/http"
	"strconv"

	"github.com/jadoint/micro/pkg/logger"
)

// Page pagination details
type Page struct {
	PageNum  int `json:"pageNum,omitempty" validate:"min=0"`
	PageSize int `json:"pageSize,omitempty" validate:"min=0"`
	Offset   int `json:"offset,omitempty" validate:"min=0"`
}

// New creates a filled instance of Page
func New(r *http.Request, pageSize int) Page {
	pageNum := GetPageNum(r)
	return Page{
		PageNum:  pageNum,
		PageSize: pageSize,
		Offset:   GetOffset(pageNum, pageSize),
	}
}

// GetPageNum gets pageNum from URL query parameter
func GetPageNum(r *http.Request) int {
	pageNumParam := r.URL.Query().Get("pageNum")
	if pageNumParam == "" {
		return 1
	}

	pageNum, err := strconv.Atoi(pageNumParam)
	if err != nil {
		logger.Log(err)
		return 1
	}
	if pageNum <= 0 {
		return 1
	}
	return pageNum
}

// GetOffset translates pageNum from a URL query
// parameter into an offset usable by a LIMIT clause.
func GetOffset(pageNum int, pageSize int) int {
	if pageNum < 1 || pageSize < 1 {
		return 0
	}
	return (pageNum - 1) * pageSize
}
