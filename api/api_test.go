package api

import (
	"io/ioutil"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/skibish/docker-demo-app/hits"
	"github.com/skibish/docker-demo-app/log"
)

func TestHostname(t *testing.T) {
	t.Parallel()

	a := &API{
		log:      log.NewNullLogger(),
		hostname: "woha",
		hits:     hits.NewHits(),
	}

	t.Run("incorrect method", func(t *testing.T) {
		ts := httptest.NewServer(a.bootRouter())
		defer ts.Close()

		resp, err := http.Post(ts.URL, "", nil)
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		if !strings.Contains(string(b), "404") {
			t.Errorf("Expected to find 404 in the response body: %v", string(b))
		}
	})

	t.Run("hostname present", func(t *testing.T) {
		ts := httptest.NewServer(a.bootRouter())
		defer ts.Close()

		resp, err := http.Get(ts.URL)
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		if !strings.Contains(string(b), "woha") {
			t.Errorf("Expected to find woha in the response body: %v", string(b))
		}
	})
}

func TestHello(t *testing.T) {
	t.Parallel()

	a := &API{
		log:      log.NewNullLogger(),
		hostname: "woha",
		hits:     hits.NewHits(),
	}

	t.Run("incorrect method", func(t *testing.T) {
		ts := httptest.NewServer(a.bootRouter())
		defer ts.Close()

		resp, err := http.Post(ts.URL+"/hello", "", nil)
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		expected := "404"
		if !strings.Contains(string(b), expected) {
			t.Errorf("Expected to find %q in the response body: %v", expected, string(b))
		}
	})

	t.Run("default hello", func(t *testing.T) {
		ts := httptest.NewServer(a.bootRouter())
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/hello")
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		expected := "Hello, world"
		if !strings.Contains(string(b), expected) {
			t.Errorf("Expected to find %q in the response body: %v", expected, string(b))
		}
	})

	t.Run("custom hello", func(t *testing.T) {
		ts := httptest.NewServer(a.bootRouter())
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/hello/test")
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		expected := "Hello, test"
		if !strings.Contains(string(b), expected) {
			t.Errorf("Expected to find %q in the response body: %v", expected, string(b))
		}
	})
}
