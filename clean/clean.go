package clean

import (
	"github.com/microcosm-cc/bluemonday"
)

// UGC configures a UGC policy then returns it.
func UGC() *bluemonday.Policy {
	ugc := bluemonday.UGCPolicy()
	ugc.AllowAttrs("style").OnElements("span", "p", "div", "em", "i", "b", "strong")
	ugc.RequireNoReferrerOnLinks(true)
	return ugc
}

// Strict returns the strictest policy that strips tags from text.
func Strict() *bluemonday.Policy {
	return bluemonday.StrictPolicy()
}
