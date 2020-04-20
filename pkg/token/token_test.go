package token_test

import (
	"testing"

	"github.com/jadoint/micro/pkg/visitor"

	"github.com/joho/godotenv"

	"github.com/jadoint/micro/pkg/token"
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
	dataClaim := visitor.GetVisitorTokenDataClaim(want.ID, want.Name)
	tokenString, err := token.Create(dataClaim)
	if err != nil {
		t.Errorf("TestMakeAuthTokenAndParseToken:MakeAuthToken failed with error: %s", err.Error())
	}
	claims, err := token.Parse(tokenString)
	if err != nil {
		t.Errorf("TestMakeAuthTokenAndParseToken:ParseToken failed with error: %s", err.Error())
	}
	got := struct {
		IAT  int64
		ID   int64
		Name string
	}{}
	got.IAT = int64(claims["iat"].(float64))
	claimsData := claims["data"].(map[string]interface{})
	got.ID = int64(claimsData["id"].(float64))
	got.Name = claimsData["name"].(string)
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
