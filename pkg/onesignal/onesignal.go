package onesignal

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jadoint/micro/pkg/logger"
)

// Config sets message details to post to OneSignal
type Config struct {
	DeviceID  string
	Title     string
	Body      string
	URL       string
	Tag       string
	SmallIcon string
	LargeIcon string
}

// Fields for OneSignal notifications
type Fields struct {
	AppID            string            `json:"app_id"`
	IncludePlayerIDs []string          `json:"include_player_ids"`
	Data             map[string]string `json:"data"`
	AndroidGroup     string            `json:"android_group"`
	Headings         map[string]string `json:"headings"`
	Contents         map[string]string `json:"contents"`
	SmallIcon        string            `json:"small_icon,omitempty"`
	LargeIcon        string            `json:"large_icon,omitempty"`
}

// NotifyDevice send notification to device through OneSignal
func NotifyDevice(cfg *Config) {
	appID := os.Getenv("ONESIGNAL_APPID")
	if appID == "" {
		logger.Log(errors.New("ONESIGNAL_APPID not set in onesignal"))
		return
	}

	fields := &Fields{
		AppID:            appID,
		IncludePlayerIDs: []string{cfg.DeviceID},
		Data:             map[string]string{"targetUrl": cfg.URL},
		AndroidGroup:     cfg.Tag,
		Headings:         map[string]string{"en": cfg.Title},
		Contents:         map[string]string{"en": cfg.Body},
		SmallIcon:        cfg.SmallIcon,
		LargeIcon:        cfg.LargeIcon,
	}

	res, err := json.Marshal(fields)
	if err != nil {
		logger.Log(errors.New("NotifyDevice: " + cfg.DeviceID))
		return
	}

	apiURL := os.Getenv("ONESIGNAL_URL")
	if apiURL == "" {
		logger.Log(errors.New("ONESIGNAL_URL not set in onesignal"))
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(res))
	if err != nil {
		logger.Log(err, "NotifyDevice NewRequest")
		return
	}
	token := os.Getenv("ONESIGNAL_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Basic "+token)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error(), "NotifyDevice - client.Do()")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error(), string(body), "NotifyDevice - ioutil.ReadAll()")
	}
}
