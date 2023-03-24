package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	todo "github.com/gopheramit/todoCLI"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()
	tempTodoFile, err := ioutil.TempFile("", "todotest")
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(newMux(tempTodoFile.Name()))
	for i := 1; i < 3; i++ {
		var body bytes.Buffer
		taskName := fmt.Sprintf("Task number%d", i)
		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}
		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}
		r, err := http.Post(ts.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}
		if r.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to add inital items:Staus:%d", r.StatusCode)
		}
	}
	return ts.URL, func() {
		ts.Close()
		os.Remove(tempTodoFile.Name())
	}

}

func TestGet(t *testing.T) {
	testCases := []struct {
		name      string
		path      string
		expCode   int
		expItems  int
		expContet string
	}{
		{name: "GetRoot", path: "/", expCode: http.StatusOK, expContet: "There is an API here"},
		{name: "GetAll", path: "/todo", expCode: http.StatusOK, expItems: 2, expContet: "Task number 1"},
		{name: "GetOne", path: "/todo/1", expCode: http.StatusOK, expItems: 1, expContet: "Task number 1"},
		{name: "NotFound", path: "/todo/500", expCode: http.StatusNotFound},
	}
	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				resp struct {
					Results      todo.List `json:"results`
					Date         int64     `json:"data"`
					TotalResults int       `json:"total_results"`
				}
				body []byte
				err  error
			)
			r, err := http.Get(url + tc.path)
			if err != nil {
				t.Error(err)
			}
			defer r.Body.Close()

			if r.StatusCode != tc.expCode {
				t.Fatalf("Expected %q,got %q", http.StatusText(tc.expCode),
					http.StatusText(r.StatusCode))
			}
			switch {
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = ioutil.ReadAll(r.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(body), tc.expContet) {
					t.Errorf("Expected %q,got %q", tc.expContet,
						string(body))
				}
			default:
				t.Fatalf("Unsupported Content-Type:%q", r.Header.Get("Content-Type"))
			}

		})
	}

}
