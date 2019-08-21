package auth_test

import (
	"testing"

	"github.com/jadoint/micro/auth"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env.testing")
}

func TestGenerateHash(t *testing.T) {
	password := "test123"
	got, err := auth.GenerateHash(password)
	if err != nil {
		t.Errorf("TestGenerateHash failed with error: %s", err.Error())
	}
	want := "$argon2id$v=19$m=65536,t=3,p=2$"
	if got[0:31] != want {
		t.Errorf("TestGenerateHash failed, got: %s, want: %s", got[0:31], want)
	}
	isMatchingPasswords, err := auth.VerifyPasswordHash(password, got)
	if err != nil {
		t.Errorf("TestGenerateHash:VerifyPasswordHash failed with error: %s", err.Error())
	}
	if !isMatchingPasswords {
		t.Errorf("TestGenerateHash:VerifyPasswordHash failed, got: %t, want: %t", isMatchingPasswords, true)
	}
	badPassword := "badpassword"
	isMatchingPasswords, err = auth.VerifyPasswordHash(badPassword, got)
	if err != nil {
		t.Errorf("TestGenerateHash:VerifyPasswordHash failed with error: %s", err.Error())
	}
	if isMatchingPasswords {
		t.Errorf("TestGenerateHash:VerifyPasswordHash failed, got: %t, want: %t", isMatchingPasswords, false)
	}
}

func TestVerifyPasswordHash(t *testing.T) {
	password := "test123"
	encodedHash := "$argon2id$v=19$m=65536,t=3,p=2$7OOSkacMICQPwQygnEwlEA$FrLtmPBc36lhfjO8QaSB0sLbHusRRsFoKOcWy5tyJsE"
	isMatchingPasswords, err := auth.VerifyPasswordHash(password, encodedHash)
	if err != nil {
		t.Errorf("TestVerifyPasswordHash failed with error: %s", err.Error())
	}
	if !isMatchingPasswords {
		t.Errorf("TestVerifyPasswordHash [good password] failed, got: %t, want: %t", isMatchingPasswords, true)
	}
	badPassword := "badpassword"
	isMatchingPasswords, err = auth.VerifyPasswordHash(badPassword, encodedHash)
	if err != nil {
		t.Errorf("TestVerifyPasswordHash failed with error: %s", err.Error())
	}
	if isMatchingPasswords {
		t.Errorf("TestVerifyPasswordHash [bad password] failed, got: %t, want: %t", isMatchingPasswords, false)
	}
}
