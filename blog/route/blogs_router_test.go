package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jadoint/micro/blog/model"
)

func TestGetRecent(t *testing.T) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		t.Skip("Set LISTEN and start server test server to run this test")
	}
	url := fmt.Sprintf("http://%s/%s/blogs/recent/1", listen, os.Getenv("START_PATH"))
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("TestGetRecent failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	// Unmarshalling
	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	got := struct {
		Blogs []*model.Blog `json:"blogs"`
	}{}
	err = d.Decode(&got)
	if err != nil {
		t.Errorf("TestGetRecent failed with error: %s", err.Error())
	}

	var prevID int64
	for _, v := range got.Blogs {
		if prevID != 0 && v.ID > prevID {
			t.Errorf("TestGetRecent failed with %d > %d", v.ID, prevID)
		}
		if len(v.Title) == 0 {
			t.Errorf("TestGetRecent failed with title length of 0")
		}
		prevID = v.ID
	}
}
