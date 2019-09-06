package msg_test

import (
	"testing"

	"github.com/jadoint/micro/pkg/msg"
)

func TestMakeAppMsg(t *testing.T) {
	appMsg := msg.MakeAppMsg("a")
	got := string(appMsg)
	want := `{"appMsg":"a"}`
	if got != want {
		t.Errorf("MakeAppMsg failed, got: %s, want: %s", got, want)
	}
}
