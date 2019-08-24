package clean

import (
	"github.com/microcosm-cc/bluemonday"
)

// UGC configures a UGC policy then returns it.
func UGC() *bluemonday.Policy {
	ugc := bluemonday.UGCPolicy()
	ugc.AllowAttrs("style").OnElements("span", "p", "div", "em", "i", "b", "strong")
	ugc.RequireNoReferrerOnLinks(true)
	// Permits the "dir", "id", "lang", "title" attributes globally
	ugc.AllowStandardAttributes()
	// Permits the "img" element and its standard attributes
	ugc.AllowImages()
	// Permits ordered and unordered lists, and also definition lists
	ugc.AllowLists()
	// Permits HTML tables and all applicable elements and non-styling attributes
	ugc.AllowTables()
	return ugc
}

// Strict returns the strictest policy that strips tags from text.
func Strict() *bluemonday.Policy {
	return bluemonday.StrictPolicy()
}
