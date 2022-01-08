package asset_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/util"
)

const (
	user, pass = "octocat", "123"
	id, file   = "Id123", "testpkg-1.2.3.tar.gz"
	location   = "http://example.org/testpkg"
)

func TestFound(t *testing.T) {
	util.TestGitHubAPI = foundAPI
	req, rec := setup()
	req.SetBasicAuth(user, pass)
	asset.ServeHTTP(rec, req)
	res := rec.Result()

	if res.StatusCode != http.StatusFound {
		t.Error("wrong status code: ", res.StatusCode)
	}
	if res.Header.Get("Location") != location {
		t.Error("wrong redirect location")
	}
}

func TestForbidden(t *testing.T) {
	util.TestGitHubAPI = foundAPI
	req, rec := setup()
	req.SetBasicAuth(user+"0", pass)

	asset.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestNotFound(t *testing.T) {
	util.TestGitHubAPI = notFoundAPI
	req, rec := setup()
	req.SetBasicAuth(user, pass)

	asset.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code: ", rec.Code)
	}
}

func setup() (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(
		http.MethodGet, path.Join(asset.BaseURLPath, id, file), nil,
	), httptest.NewRecorder()
}

func foundAPI(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, `{"data": {
		"node": {
			"url": "%s",
			"release": { "repository": { "owner": {
				"login": "%s"
			}}}
		}
	}}`, location, user)
}

func notFoundAPI(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, `{
		"data": {
			"node": null
		},
		"errors": [{
			"type": "NOT_FOUND",
			"path": [ "node" ],
			"locations": [{ "line": 2, "column": 3 }],
			"message": "Could not resolve to a node with the global id of '%s'"
		}]
	}`, id)
}
