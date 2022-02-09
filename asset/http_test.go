package asset

import (
	"net/http"
	"testing"
)

const (
	user, pass = "octocat", "123"
	id, file   = "Id123", "octopack-1.2.3.tar.gz"
	location   = "http://example.org/octopack"
)

func TestFound(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	wraph(makeFoundAPI(t)).ServeHTTP(rec, req)
	res := rec.Result()

	if res.StatusCode != http.StatusFound {
		t.Error("wrong status code:", res.StatusCode)
	}
	if res.Header.Get("Location") != location {
		t.Error("wrong redirect location")
	}
}

func TestForbidden(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user+"0", pass)

	wraph(makeFoundAPI(t)).ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Error("wrong status code:", rec.Code)
	}
}

func TestNotFound(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	wraph(notFoundAPI).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code:", rec.Code)
	}
}

func TestUnauth(t *testing.T) {
	req, rec := setup()
	h := wraph(func(rw http.ResponseWriter, r *http.Request) {
		t.Error("should not invoke GitHub API")
		http.NotFound(rw, r)
	})

	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Error("wrong status code:", rec.Code)
	}
}
