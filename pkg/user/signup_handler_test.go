package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/db"
	"github.com/jadoint/micro/pkg/validate"
)

func TestSignupSuccess(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/signup", listen, os.Getenv("START_PATH"))
	username := "TestSignupUser"
	password := "test123"
	email := "test@gmail.com"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestSignupSuccess failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	// Unmarshalling
	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	newUser := struct {
		ID       int64  `json:"id" validate:"required,min=1"`
		Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	}{}

	err = d.Decode(&newUser)
	if err != nil {
		t.Errorf("TestSignupSuccess failed with error: %s", err.Error())
	}

	// Validation
	err = validate.Struct(newUser)
	if err != nil {
		t.Errorf("TestSignupSuccess failed with error: %s", err.Error())
	}

	if newUser.ID == 0 {
		t.Errorf("TestSignupSuccess failed, got: %d, want %s", newUser.ID, "> 0")
	}

	if newUser.Username != username {
		t.Errorf("TestSignupSuccess failed, got: %s, want: %s", newUser.Username, username)
	}

	// Verify the new user is in the database
	// Database
	dbClient, err := db.GetClient()
	if err != nil {
		t.Errorf("TestSignupSuccess failed with error: %s", err.Error())
	}

	// Clients
	clients := &conn.Clients{DB: dbClient}
	defer clients.DB.Master.Close()
	defer clients.DB.Read.Close()

	dbUser, err := GetUser(clients, newUser.ID)
	if err != nil {
		t.Errorf("TestSignupSuccess failed with error: %s", err.Error())
	}
	if dbUser.ID != newUser.ID {
		t.Errorf("TestSignupSuccess failed, got: %d, want %d", dbUser.ID, newUser.ID)
	}
	if dbUser.Username != newUser.Username {
		t.Errorf("TestSignupSuccess failed, got: %s, want %s", dbUser.Username, newUser.Username)
	}

	// DB cleanup
	_, err = clients.DB.Exec(`
		DELETE FROM user
		WHERE id_user = ?
		LIMIT 1`,
		dbUser.ID)
	if err != nil {
		t.Errorf("TestSignupSuccess:Cleanup failed with error: %s", err.Error())
	}
}

func TestSignupBadUsername(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/signup", listen, os.Getenv("START_PATH"))
	username := "!@#$%^&*())_+"
	password := "test123"
	email := "test@gmail.com"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestSignupBadUsername failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestSignupBadUsername failed with error: %s", err.Error())
	}

	got := string(body)
	want := `"error":`
	if !strings.Contains(got, want) {
		t.Errorf("TestSignupBadUsername failed, got: %s, want: %s", got, want)
	}
}

func TestSignupBadPassword(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/signup", listen, os.Getenv("START_PATH"))
	username := "TestSignupBadPassword"
	password := "123"
	email := "test@gmail.com"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestSignupBadPassword failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestSignupBadPassword failed with error: %s", err.Error())
	}

	got := string(body)
	want := `"error":`
	if !strings.Contains(got, want) {
		t.Errorf("TestSignupBadPassword failed, got: %s, want: %s", got, want)
	}
}

func TestSignupBadEmail(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/signup", listen, os.Getenv("START_PATH"))
	username := "TestSignupBadEmail"
	password := "test123"
	email := "bademail"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestSignupBadEmail failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestSignupBadEmail failed with error: %s", err.Error())
	}

	got := string(body)
	want := `"error":`
	if !strings.Contains(got, want) {
		t.Errorf("TestSignupBadEmail failed, got: %s, want: %s", got, want)
	}
}

func TestSignupAlreadyLoggedIn(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/signup", listen, os.Getenv("START_PATH"))
	username := "TestSignupUser"
	password := "test123"
	email := "test@gmail.com"
	postFields := fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestSignupAlreadyLoggedIn failed with error: %s", err.Error())
	}

	req.AddCookie(&http.Cookie{Name: os.Getenv("COOKIE_SESSION_NAME"), Value: os.Getenv("TEST_USER_COOKIE")})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("TestSignupAlreadyLoggedIn failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestSignupAlreadyLoggedIn failed with error: %s", err.Error())
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("TestSignupAlreadyLoggedIn failed, got: %d, want: %d", resp.StatusCode, http.StatusOK)
	}

	got := string(body)
	want := `"error":`
	if !strings.Contains(got, want) {
		t.Errorf("TestSignupAlreadyLoggedIn failed, got: %s, want: %s", got, want)
	}
}
