package simple

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/util"
)

const (
	user, pass = "octocat", "123"
	pkg, repo  = "octopack", "test-repo"
)

var assets = []ghAsset{
	{"Id0", pkg + "-1.1.1.tar.gz"},
	{"Id1", pkg + "-1.2.0.tar.gz"},
	{"Id2", pkg + "-1.2.0-py3-none-any.whl"},
}

func TestRoot(t *testing.T) {
	req, rec := setup("")
	req.SetBasicAuth(user, pass)

	wraph(makeUnreachableAPI(t)).ServeHTTP(rec, req)
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
	req, rec := setup(pkg)
	req.SetBasicAuth(user, pass)

	wraph(makeFoundAPI(t)).ServeHTTP(rec, req)
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
			asset.MakeURL(a.ID, a.Name), a.Name,
		)
		if !strings.Contains(body, anchor) {
			t.Error("missing link:", anchor)
		}
	}
}

func TestNotFound(t *testing.T) {
	req, rec := setup(pkg)
	req.SetBasicAuth(user, pass)

	wraph(notFoundAPI).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestUnauth(t *testing.T) {
	req, rec := setup(pkg)

	wraph(makeUnreachableAPI(t)).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Error("wrong status code: ", rec.Code)
	}
}

func setup(pkg string) (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(
		http.MethodGet, path.Join(pathBase, pkg)+"/", nil,
	), httptest.NewRecorder()
}

func wraph(serve http.HandlerFunc) http.Handler {
	return &handler{util.NewGHv4ClientMaker(serve)}
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
			assets[0].ID, assets[0].Name,
			assets[1].ID, assets[1].Name,
			assets[2].ID, assets[2].Name,
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
