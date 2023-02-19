package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jadoint/micro/pkg/validate"
)

func TestLoginSuccess(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/login", listen, os.Getenv("START_PATH"))
	username := "testuser1"
	password := "test123"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestLoginSuccess failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	// Unmarshalling
	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	newUser := struct {
		ID       int    `json:"id" validate:"required,min=1"`
		Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	}{}

	err = d.Decode(&newUser)
	if err != nil {
		t.Errorf("TestLoginSuccess failed with error: %s", err.Error())
	}

	// Validation
	err = validate.Struct(newUser)
	if err != nil {
		t.Errorf("TestLoginSuccess failed with error: %s", err.Error())
	}

	if newUser.ID == 0 {
		t.Errorf("TestLoginSuccess failed, got: %d, want %s", newUser.ID, "> 0")
	}

	if newUser.Username != username {
		t.Errorf("TestLoginSuccess failed, got: %s, want: %s", newUser.Username, username)
	}
}

func TestLoginBadPassword(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/login", listen, os.Getenv("START_PATH"))
	username := "testuser1"
	password := "badpassword"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestLoginBadPassword failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestLoginBadPassword failed with error: %s", err.Error())
	}

	got := string(body)
	want := `"error":`
	if !strings.Contains(got, want) {
		t.Errorf("TestLoginBadPassword failed, got: %s, want: %s", got, want)
	}
}

func TestLoginBadUsername(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/login", listen, os.Getenv("START_PATH"))
	username := "badusername"
	password := "test123"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestLoginBadUsername failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestLoginBadUsername failed with error: %s", err.Error())
	}

	got := string(body)
	want := `"error":`
	if !strings.Contains(got, want) {
		t.Errorf("TestLoginBadUsername failed, got: %s, want: %s", got, want)
	}
}

func TestLoginAlreadyLoggedIn(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/login", listen, os.Getenv("START_PATH"))
	username := "testuser1"
	password := "test123"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestLoginAlreadyLoggedIn failed with error: %s", err.Error())
	}

	req.AddCookie(&http.Cookie{Name: os.Getenv("COOKIE_SESSION_NAME"), Value: os.Getenv("TEST_USER_COOKIE")})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("TestLoginAlreadyLoggedIn failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestLoginAlreadyLoggedIn failed with error: %s", err.Error())
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("TestLoginAlreadyLoggedIn failed, got: %d, want: %d", resp.StatusCode, http.StatusOK)
	}

	got := string(body)
	want := `"error":`
	if !strings.Contains(got, want) {
		t.Errorf("TestLoginAlreadyLoggedIn failed, got: %s, want: %s", got, want)
	}
}
