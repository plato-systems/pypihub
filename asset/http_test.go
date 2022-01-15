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

func TestFound(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	makeFound(t).ServeHTTP(rec, req)
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

	makeFound(t).ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestNotFound(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	h := handler{mockClient{t: t, err: errors.New("not found")}}
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestUnauth(t *testing.T) {
	req, rec := setup()

	h := handler{mockClient{t: t, noreach: true}}
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

func makeFound(t *testing.T) *handler {
	mc := mockClient{t: t, a: ghAsset{URL: location}}
	mc.a.Release.Repository.Owner.Login = user
	return &handler{mc}
}

type mockClient struct {
	noreach bool
	err     error
	a       ghAsset
	t       *testing.T
}

func (m mockClient) Query(
	ctx context.Context, token string,
	q interface{}, v map[string]interface{},
) error {
	if m.noreach {
		m.t.Error("should not call API")
		return m.err
	}

	query, ok := q.(*queryAsset)
	if !ok {
		m.t.Error("incorrect query type")
		return m.err
	}

	query.Node.ReleaseAsset = m.a
	return m.err
}
