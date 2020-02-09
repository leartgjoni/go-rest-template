package http

import (
	mock "github.com/leartgjoni/go-rest-template/mock/http"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestServerRoutes(t *testing.T) {
	var tests = []struct {
		method          string
		route           string
		expectedInvoked []string
	}{
		{
			"POST",
			"/auth/signup",
			[]string{"AuthHandler.HandleSignup"},
		},
		{
			"POST",
			"/auth/login",
			[]string{"AuthHandler.HandleLogin"},
		},
		{
			"GET",
			"/auth/me",
			[]string{"AuthHandler.Authentication", "AuthHandler.HandleMe"},
		},
		{
			"GET",
			"/articles",
			[]string{"ArticleHandler.HandleList"},
		},
		{
			"GET",
			"/articles/random-slug",
			[]string{"ArticleHandler.ArticleCtx", "ArticleHandler.HandleGet"},
		},
		{
			"POST",
			"/articles",
			[]string{"AuthHandler.Authentication", "ArticleHandler.HandleCreate"},
		},
		{
			"PATCH",
			"/articles/random-slug",
			[]string{"AuthHandler.Authentication", "ArticleHandler.ArticleCtx", "ArticleHandler.ArticleOwner", "ArticleHandler.HandleUpdate"},
		},
		{
			"DELETE",
			"/articles/random-slug",
			[]string{"AuthHandler.Authentication", "ArticleHandler.ArticleCtx", "ArticleHandler.ArticleOwner", "ArticleHandler.HandleDelete"},
		},
	}

	for _, test := range tests {
		server := NewServer()
		invoked := &[]string{}
		// mock handlers
		server.articleHandler = mock.NewMockArticleHandler(invoked)
		server.authHandler = mock.NewMockAuthHandler(invoked)

		router := server.router()

		w := httptest.NewRecorder()
		r, err := http.NewRequest(test.method, test.route, nil)
		if err != nil {
			t.Fatal("creating request failed", err)
		}

		router.ServeHTTP(w, r)

		if !reflect.DeepEqual(*invoked, test.expectedInvoked) {
			t.Errorf("Expect %s but got %v", test.expectedInvoked, *invoked)
		}
	}
}
