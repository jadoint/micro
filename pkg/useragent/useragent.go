package useragent

import (
	"net/http"
	"strings"
)

// IsBot determines if a request is from a bot
func IsBot(r *http.Request) bool {
	bots := []string{"google", "bingbot", "yahoo"}
	ua := r.Header.Get("User-Agent")
	for _, bot := range bots {
		if strings.Contains(ua, bot) {
			return true
		}
	}
	return false
}
