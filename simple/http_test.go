package simple_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/simple"
	"github.com/plato-systems/pypihub/util"
)

const (
	user, pass = "octocat", "123"
	pkg, repo  = "octopack", "test-repo"
)

func TestRoot(t *testing.T) {
	util.TestGitHubAPI = makeUnreachableAPI(t)
	req, rec := setup("")
	req.SetBasicAuth(user, pass)

	simple.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestFoundPkg(t *testing.T) {
	testFound(t, pkg)
}

func TestFoundRepo(t *testing.T) {
	testFound(t, repo)
}

func testFound(t *testing.T, pkg string) {
	util.LoadConfigFile("./testdata/octo.toml")
	util.TestGitHubAPI = makeFoundAPI(t)
	req, rec := setup(pkg)
	req.SetBasicAuth(user, pass)

	simple.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Error("wrong status code: ", rec.Code)
	}

	bodyBuf, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatal("could not read response body: ", err)
	}
	body := string(bodyBuf)
	for _, a := range assets {
		anchor := fmt.Sprintf(
			`<a href="%s">%s</a>`,
			asset.MakeURL(a.id, a.name), a.name,
		)
		if !strings.Contains(body, anchor) {
			t.Error("missing link:", anchor)
		}
	}
}

func TestNotFound(t *testing.T) {
	util.TestGitHubAPI = notFoundAPI
	req, rec := setup(pkg)
	req.SetBasicAuth(user, pass)

	simple.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestUnauth(t *testing.T) {
	util.TestGitHubAPI = makeUnreachableAPI(t)
	req, rec := setup(pkg)

	simple.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Error("wrong status code: ", rec.Code)
	}
}

type ghAsset struct {
	id, name string
}

var assets = []ghAsset{
	{"Id0", pkg + "-1.1.1.tar.gz"},
	{"Id1", pkg + "-1.2.0.tar.gz"},
	{"Id2", pkg + "-1.2.0-py3-none-any.whl"},
}

func setup(pkg string) (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(
		http.MethodGet, path.Join(simple.BaseURLPath, pkg)+"/", nil,
	), httptest.NewRecorder()
}

func makeUnreachableAPI(t *testing.T) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t.Error("should not invoke GitHub API")
		http.NotFound(rw, r)
	}
}

func makeFoundAPI(t *testing.T) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bodyBuf, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error("could not read GraphQL query body: ", err)
		} else {
			body := string(bodyBuf)
			if strings.Contains(body, pkg) || !strings.Contains(body, repo) {
				t.Error("Package not converted to Repo in query: ", body)
			}
		}
		fmt.Fprintf(
			rw, `{ "data": { "repository": { "releases": {
				"nodes": [
					{ "releaseAssets": { "nodes": [] } },
					{ "releaseAssets": { "nodes": [{
						"id": "%s", "name": "%s"
					}]}},
					{ "releaseAssets": { "nodes": [
						{ "id": "%s", "name": "%s" },
						{ "id": "%s", "name": "%s" }
					]}}
				],
				"pageInfo": { "endCursor": "c0", "hasNextPage": false }
			}}}}`,
			assets[0].id, assets[0].name,
			assets[1].id, assets[1].name,
			assets[2].id, assets[2].name,
		)
	}
}

func notFoundAPI(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, `{
		"data": { "repository": null },
		"errors": [{
			"type": "NOT_FOUND",
			"path": [ "repository" ],
			"locations": [{ "line": 17, "column": 3 }],
			"message": "Could not resolve to a Repository with the name '%s/%s'."
		}]
	}`, user, repo)
}
