package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jadoint/micro/now"

	"github.com/jadoint/micro/blog/model"
	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/db"
	"github.com/jadoint/micro/validate"
)

func TestGetBlog(t *testing.T) {
	url := fmt.Sprintf("http://%s/%s/blog/1", os.Getenv("LISTEN"), os.Getenv("START_PATH"))
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("TestGetBlog failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	// Unmarshalling
	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	var got model.Blog
	err = d.Decode(&got)
	if err != nil {
		t.Errorf("TestGetBlog failed with error: %s", err.Error())
	}

	want := &model.Blog{
		ID:        1,
		IDAuthor:  1,
		Title:     "Lorem Ipsum Dolor 1",
		Post:      "<p>Lorem Ipsum Dolor Sit Amet</p>",
		WordCount: 5,
		Created:   "2019-08-20 09:07:14",
		Modified:  "2019-08-20 09:07:14",
	}

	if got.ID != want.ID {
		t.Errorf("TestGetBlog failed, got: %d, want: %d", got.ID, want.ID)
	}
	if got.IDAuthor != want.IDAuthor {
		t.Errorf("TestGetBlog failed, got: %d, want: %d", got.IDAuthor, want.IDAuthor)
	}
	if got.Title != want.Title {
		t.Errorf("TestGetBlog failed, got: %s, want: %s", got.Title, want.Title)
	}
	if got.Post != want.Post {
		t.Errorf("TestGetBlog failed, got: %s, want: %s", got.Post, want.Post)
	}
	if got.WordCount != want.WordCount {
		t.Errorf("TestGetBlog failed, got: %d, want: %d", got.WordCount, want.WordCount)
	}
	if got.Created != want.Created {
		t.Errorf("TestGetBlog failed, got: %s, want: %s", got.Created, want.Created)
	}
	if got.Modified != want.Modified {
		t.Errorf("TestGetBlog failed, got: %s, want: %s", got.Modified, want.Modified)
	}
}

func TestGetBlogWithVisitor(t *testing.T) {
	url := fmt.Sprintf("http://%s/%s/blog/1", os.Getenv("LISTEN"), os.Getenv("START_PATH"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{Name: os.Getenv("COOKIE_SESSION_NAME"), Value: os.Getenv("TEST_USER_COOKIE")})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestGetBlogWithVisitor failed with error: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestGetBlogWithVisitor failed, got: %d, want: %d", resp.StatusCode, http.StatusOK)
	}

	table := map[string]interface{}{
		"idPost":    1,
		"idAuthor":  1,
		"title":     `"Lorem Ipsum Dolor 1"`,
		"idVisitor": 1,
	}

	got := string(body)
	for key, val := range table {
		want := fmt.Sprintf(`"%s":%v`, key, val)
		if !strings.Contains(got, want) {
			t.Errorf("TestGetBlogWithVisitor failed, got: %s, want: %s", got, want)
		}
	}
}

func TestPostBlogSuccess(t *testing.T) {
	url := fmt.Sprintf("http://%s/%s/blog", os.Getenv("LISTEN"), os.Getenv("START_PATH"))
	test := &model.Blog{
		Title:      "<h1>Test Title</h1>",
		Post:       "Test Post<script>alert('test')</script>",
		IsUnlisted: false,
		IsDraft:    false,
	}
	postFields := fmt.Sprintf(`{"title": "%s", "post": "%s", "isUnlisted": %t, "isDraft": %t}`, test.Title, test.Post, test.IsUnlisted, test.IsDraft)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestPostBlogSuccess failed with error: %s", err.Error())
	}

	req.AddCookie(&http.Cookie{Name: os.Getenv("COOKIE_SESSION_NAME"), Value: os.Getenv("TEST_USER_COOKIE")})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("TestPostBlogSuccess failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestPostBlogSuccess failed, got: %d, want: %d", resp.StatusCode, http.StatusOK)
	}

	// Unmarshalling
	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	newBlog := struct {
		ID int64 `json:"idPost" validate:"required,min=1"`
	}{}

	err = d.Decode(&newBlog)
	if err != nil {
		t.Errorf("TestPostBlogSuccess failed with error: %s", err.Error())
	}

	// Validation
	err = validate.Struct(newBlog)
	if err != nil {
		t.Errorf("TestPostBlogSuccess failed with error: %s", err.Error())
	}

	if newBlog.ID == 0 {
		t.Errorf("TestPostBlogSuccess failed, got: %d, want %s", newBlog.ID, "> 0")
	}

	// Verify the new blog is in the database
	// Database
	dbClient, err := db.GetClient()
	if err != nil {
		t.Errorf("TestPostBlogSuccess failed with error: %s", err.Error())
	}

	// Clients
	clients := &conn.Clients{DB: dbClient}
	defer clients.DB.Master.Close()
	defer clients.DB.Read.Close()

	dbBlog, err := model.GetBlog(clients, newBlog.ID)
	if err != nil {
		t.Errorf("TestPostBlogSuccess failed with error: %s", err.Error())
	}

	if dbBlog.ID != newBlog.ID {
		t.Errorf("TestPostBlogSuccess failed, got: %d, want %d", dbBlog.ID, newBlog.ID)
	}

	want := &model.Blog{
		Title:      "Test Title",
		Post:       "Test Post",
		IsUnlisted: false,
		IsDraft:    false,
	}

	if dbBlog.Title != want.Title {
		t.Errorf("TestPostBlogSuccess failed, got: %s, want %s", dbBlog.Title, want.Title)
	}
	if dbBlog.Post != want.Post {
		t.Errorf("TestPostBlogSuccess failed, got: %s, want %s", dbBlog.Post, want.Post)
	}

	// DB cleanup
	err = model.DeleteBlog(clients, newBlog.ID)
	if err != nil {
		t.Errorf("TestPostBlogSuccess:Cleanup failed with error: %s", err.Error())
	}
}

func TestUpdateBlogSuccess(t *testing.T) {
	url := fmt.Sprintf("http://%s/%s/blog/6", os.Getenv("LISTEN"), os.Getenv("START_PATH"))
	test := &model.Blog{
		ID:         6,
		Title:      "<h1>Updated</h1>",
		Post:       "New Update<script>alert('test')</script>",
		WordCount:  2,
		IsUnlisted: true,
		IsDraft:    true,
		Modified:   now.MySQLUTC(),
	}
	postFields := fmt.Sprintf(`{"title": "%s", "post": "%s", "isUnlisted": %t, "isDraft": %t}`, test.Title, test.Post, test.IsUnlisted, test.IsDraft)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(postFields)))
	if err != nil {
		t.Errorf("TestPostBlogSuccess failed with error: %s", err.Error())
	}

	req.AddCookie(&http.Cookie{Name: os.Getenv("COOKIE_SESSION_NAME"), Value: os.Getenv("TEST_USER_COOKIE")})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("TestUpdateBlogSuccess failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestUpdateBlogSuccess failed, got: %d, want: %d", resp.StatusCode, http.StatusOK)
	}

	// Unmarshalling
	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	b := struct {
		ID int64 `json:"idPost" validate:"required,min=1"`
	}{}

	err = d.Decode(&b)
	if err != nil {
		t.Errorf("TestUpdateBlogSuccess failed with error: %s", err.Error())
	}

	// Validation
	err = validate.Struct(b)
	if err != nil {
		t.Errorf("TestUpdateBlogSuccess failed with error: %s", err.Error())
	}

	want := &model.Blog{
		ID:         6,
		Title:      "Updated",
		Post:       "New Update",
		WordCount:  2,
		IsUnlisted: true,
		IsDraft:    true,
		Modified:   now.MySQLUTC(),
	}

	if b.ID != test.ID {
		t.Errorf("TestUpdateBlogSuccess failed, got: %d, want %d", b.ID, want.ID)
	}

	// Verify the update in the database
	// Database
	dbClient, err := db.GetClient()
	if err != nil {
		t.Errorf("TestUpdateBlogSuccess failed with error: %s", err.Error())
	}

	// Clients
	clients := &conn.Clients{DB: dbClient}
	defer clients.DB.Master.Close()
	defer clients.DB.Read.Close()

	got, err := model.GetBlog(clients, b.ID)
	if err != nil {
		t.Errorf("TestUpdateBlogSuccess failed with error: %s", err.Error())
	}

	if got.ID != want.ID {
		t.Errorf("TestUpdateBlogSuccess failed, got: %d, want %d", got.ID, want.ID)
	}
	if got.Title != want.Title {
		t.Errorf("TestUpdateBlogSuccess failed, got: %s, want %s", got.Title, want.Title)
	}
	if got.Post != want.Post {
		t.Errorf("TestUpdateBlogSuccess failed, got: %s, want %s", got.Post, want.Post)
	}
	// if got.IsUnlisted != want.IsUnlisted {
	// 	t.Errorf("TestUpdateBlogSuccess failed, got: %t, want %t", got.IsUnlisted, want.IsUnlisted)
	// }
	// if got.IsDraft != want.IsDraft {
	// 	t.Errorf("TestUpdateBlogSuccess failed, got: %t, want %t", got.IsDraft, want.IsDraft)
	// }
	if got.Modified < want.Modified {
		t.Errorf("TestUpdateBlogSuccess failed, got: %s, want %s", got.Modified, want.Modified)
	}

	// DB cleanup - revert all fields to original except for `modified`
	original := &model.Blog{
		ID:         6,
		Title:      "Test Update",
		Post:       "<p>Lorem Ipsum Dolor Sit Amet</p>",
		WordCount:  5,
		IsUnlisted: false,
		IsDraft:    false,
		Modified:   want.Modified,
	}
	err = model.UpdateBlog(clients, original)
	if err != nil {
		t.Errorf("TestUpdateBlogSuccess:Cleanup failed with error: %s", err.Error())
	}
}

func TestDeleteBlogSuccess(t *testing.T) {
	// DB setup
	dbClient, err := db.GetClient()
	if err != nil {
		t.Errorf("TestUpdateBlogSuccess failed with error: %s", err.Error())
	}

	// Clients
	clients := &conn.Clients{DB: dbClient}
	defer clients.DB.Master.Close()
	defer clients.DB.Read.Close()

	b := &model.Blog{
		IDAuthor:   1,
		Title:      "Test Delete",
		Post:       "Delete",
		WordCount:  1,
		IsUnlisted: false,
		IsDraft:    false,
	}
	idBlog, err := model.AddBlog(clients, b)
	if err != nil {
		t.Errorf("TestDeleteBlogSuccess:Setup failed with error: %s", err.Error())
	}

	// Check inserted blog
	dbBlog, err := model.GetBlog(clients, idBlog)
	if err != nil {
		t.Errorf("TestDeleteBlogSuccess failed with error: %s", err.Error())
	}
	if dbBlog.ID != idBlog {
		t.Errorf("TestDeleteBlogSuccess failed, got: %d, want %d", dbBlog.ID, idBlog)
	}

	// Finally testing Delete API
	url := fmt.Sprintf("http://%s/%s/blog/%d", os.Getenv("LISTEN"), os.Getenv("START_PATH"), idBlog)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Errorf("TestDeleteBlogSuccess failed with error: %s", err.Error())
	}

	req.AddCookie(&http.Cookie{Name: os.Getenv("COOKIE_SESSION_NAME"), Value: os.Getenv("TEST_USER_COOKIE")})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("TestDeleteBlogSuccess failed with error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestDeleteBlogSuccess failed, got: %d, want: %d", resp.StatusCode, http.StatusOK)
	}

	// Unmarshalling
	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	rb := struct {
		ID int64 `json:"idPost" validate:"required,min=1"`
	}{}

	err = d.Decode(&rb)
	if err != nil {
		t.Errorf("TestDeleteBlogSuccess failed with error: %s", err.Error())
	}

	// Validation
	err = validate.Struct(rb)
	if err != nil {
		t.Errorf("TestDeleteBlogSuccess failed with error: %s", err.Error())
	}

	if rb.ID == 0 {
		t.Errorf("TestDeleteBlogSuccess failed, got: %d, want %s", rb.ID, "> 0")
	}

	// Verify the delete in the database
	_, err = model.GetBlog(clients, idBlog)
	if err == nil {
		t.Errorf("TestDeleteBlogSuccess failed (deleted entry found) with error: %s", err.Error())
	}
}
