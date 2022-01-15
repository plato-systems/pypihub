package asset

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	user, pass = "octocat", "123"
	id, file   = "Id123", "octopack-1.2.3.tar.gz"
	location   = "http://example.org/octopack"
)

type mockAPI struct {
	noreach *testing.T
	err     error
	a       asset
}

func (m mockAPI) getAsset(ctx context.Context, token, id string) (asset, error) {
	if m.noreach != nil {
		m.noreach.Error("should not call API")
	}
	return m.a, m.err
}

var found = handler{mockAPI{a: asset{url: location, owner: user}}}

func TestFound(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	found.ServeHTTP(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusFound {
		t.Error("wrong status code: ", res.StatusCode)
	}
	if res.Header.Get("Location") != location {
		t.Error("wrong redirect location")
	}
}

func TestForbidden(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user+"0", pass)

	found.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestNotFound(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	h := handler{mockAPI{err: errors.New("not found")}}
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestUnauth(t *testing.T) {
	req, rec := setup()

	h := handler{mockAPI{noreach: t}}
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Error("wrong status code: ", rec.Code)
	}
}

func setup() (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(
		http.MethodGet, MakeURL(id, file), nil,
	), httptest.NewRecorder()
}
