package asset2_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	asset "github.com/plato-systems/pypihub/asset2"
	"github.com/plato-systems/pypihub/util"
)

const (
	user, pass = "octocat", "123"
	id, file   = "Id123", "octopack-1.2.3.tar.gz"
	location   = "http://example.org/octopack"
)

type mockAPI struct {
	noreach *testing.T
	*asset.Asset
}

var foundAPI = mockAPI{Asset: &asset.Asset{URL: location, Owner: user}}

func (m mockAPI) GetAsset(ctx context.Context, token, id string) (asset.Asset, error) {
	a := asset.Asset{}
	if m.noreach != nil {
		m.noreach.Error("should not call API")
		return a, nil
	}
	if m.Asset == nil {
		return a, fmt.Errorf("no such asset")
	}
	return *m.Asset, nil
}

func TestFound(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	asset.NewHandler(foundAPI).ServeHTTP(rec, req)
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

	asset.NewHandler(foundAPI).ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestNotFound(t *testing.T) {
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	asset.NewHandler(mockAPI{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestUnauth(t *testing.T) {
	util.TestGitHubAPI = func(rw http.ResponseWriter, r *http.Request) {
		t.Error("should not invoke GitHub API")
		http.NotFound(rw, r)
	}
	req, rec := setup()

	asset.NewHandler(mockAPI{noreach: t}).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Error("wrong status code: ", rec.Code)
	}
}

func setup() (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(
		http.MethodGet, asset.MakeURL(id, file), nil,
	), httptest.NewRecorder()
}
