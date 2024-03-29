package user

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestLogoutSuccess(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/auth/logout", listen, os.Getenv("START_PATH"))
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Errorf("TestLogoutSuccess failed with error: %s", err.Error())
	}

	req.AddCookie(&http.Cookie{Name: os.Getenv("COOKIE_SESSION_NAME"), Value: os.Getenv("TEST_USER_COOKIE")})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestLogoutSuccess failed, got: %d, want: %d", resp.StatusCode, http.StatusOK)
	}

	got := string(body)
	want := `"appMsg":`
	if !strings.Contains(got, want) {
		t.Errorf("TestLogoutSuccess failed, got: %s, want: %s", got, want)
	}
}
