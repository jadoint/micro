package msg

import (
	"encoding/json"

	"github.com/jadoint/micro/pkg/logger"
)

// AppMessage message to be sent by the app to the user
type AppMessage struct {
	AppMsg string `json:"appMsg"`
}

// MakeAppMsg builds message to be sent by the app to the user
func MakeAppMsg(msg string) []byte {
	am := &AppMessage{AppMsg: msg}
	res, err := json.Marshal(am)
	if err != nil {
		logger.Log(err)
		return nil
	}
	return res
}
