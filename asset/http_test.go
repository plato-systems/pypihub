package asset_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/util"
)

const (
	user, pass = "octocat", "123"
	id, file   = "Id123", "octopack-1.2.3.tar.gz"
	location   = "http://example.org/octopack"
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

func TestUnauth(t *testing.T) {
	util.TestGitHubAPI = func(rw http.ResponseWriter, r *http.Request) {
		t.Error("should not invoke GitHub API")
		http.NotFound(rw, r)
	}
	req, rec := setup()

	asset.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Error("wrong status code: ", rec.Code)
	}
}

func setup() (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(
		http.MethodGet, asset.MakeURL(id, file), nil,
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
