package http

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer()
	server.Addr = ":1234"
	if err := server.Open(); err != nil {
		t.Fatal("Error on server.Open()", err)
	}
	defer server.Close()

	resp, err := http.Get("http://localhost:8181/health")
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
