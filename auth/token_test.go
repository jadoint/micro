package auth_test

import (
	"testing"

	"github.com/joho/godotenv"

	"github.com/jadoint/micro/auth"
)

func init() {
	godotenv.Load(".env.testing")
}

func TestMakeAuthTokenAndParseToken(t *testing.T) {
	want := struct {
		ID   int64
		Name string
	}{
		ID:   1,
		Name: "Username",
	}
	shortToken, err := auth.MakeAuthToken(want.ID, want.Name)
	if err != nil {
		t.Errorf("TestMakeAuthTokenAndParseToken:MakeAuthToken failed with error: %s", err.Error())
	}
	got, err := auth.ParseToken(shortToken)
	if err != nil {
		t.Errorf("TestMakeAuthTokenAndParseToken:ParseToken failed with error: %s", err.Error())
	}
	if got.IAT == 0 {
		t.Errorf("TestMakeAuthTokenAndParseToken failed, got.ID: %d, want: %s", got.IAT, "> 0")
	}
	if got.ID != want.ID {
		t.Errorf("TestMakeAuthTokenAndParseToken failed, got.ID: %d, want: %d", got.ID, want.ID)
	}
	if got.Name != want.Name {
		t.Errorf("TestMakeAuthTokenAndParseToken failed, got.ID: %s, want: %s", got.Name, want.Name)
	}
}
