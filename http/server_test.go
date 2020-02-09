package http

import (
	"fmt"
	mock "github.com/leartgjoni/go-rest-template/mock/http"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestServerListeningIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	server := NewServer()
	server.Addr = ":1234"
	if err := server.Start(); err != nil {
		t.Fatal("Error on server.Open()", err)
	}
	defer func() {
		if err := server.Close(); err != nil {
			t.Errorf("error closing server: %s", err)
		}
	}()

	resp, err := http.Get("http://localhost:1234/health")
	if err != nil {
		t.Fatal("http get failed", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("http get failed", err)
	}

	if string(body) != "healthy" {
		t.Fatalf("Expected 'healthy' but got %s", string(body))
	}
}

func TestServerRoutesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	server := NewServer()
	// mock handlers
	articleHandler := &mock.ArticleHandler{}
	authHandler := &mock.AuthHandler{}
	server.articleHandler = articleHandler
	server.authHandler = authHandler

	server.Addr = ":1235"

	if err := server.Open(); err != nil {
		t.Fatal("Error on server.Open()", err)
	}
	defer server.Close()

	var tests = []struct {
		method string
		route string
		expectedInvoked []string
	}{
		{
			"GET",
			"articles",
			[]string{"HandleList"},
		},
		{
				"GET",
			"articles/random-slug",
			[]string{"ArticleCtx", "HandleGet"},
		},
	}

	for _, test := range tests {
		client := http.Client{}
		req, err := http.NewRequest(test.method, fmt.Sprintf("http://localhost:1235/%s", test.route), nil)
		if err != nil {
			t.Fatal("creating request failed", err)
		}

		if _, err := client.Do(req); err != nil {
			t.Fatal("http call failed", err)
		}

		// TODO: solve problem when Invoked come from two different handlers. Like on save article (auth and article in this case)
		if !reflect.DeepEqual(articleHandler.Invoked, test.expectedInvoked) {
			t.Errorf("Expect %s but got %s", test.expectedInvoked, articleHandler.Invoked)
		}

		// reset mocks
		authHandler.Invoked = []string{}
		articleHandler.Invoked = []string{}
	}
}
