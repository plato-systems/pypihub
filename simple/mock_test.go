package simple

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/plato-systems/pypihub/util"
)

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
		i := verifyBody(t, r.Body)
		next := fmt.Sprintf("%s%d", cursor, i+1)
		fmt.Fprintf(
			rw, `{ "data": { "repository": { "releases": {
				"nodes": [
					{ "releaseAssets": { "nodes": [] } },
					{ "releaseAssets": { "nodes": [
						{ "id": "%s", "name": "%s" }
					]}}
				],
				"pageInfo": { "endCursor": "%s", "hasNextPage": %t }
			}}}}`,
			assets[i].ID, assets[i].Name, next, i+1 < len(assets),
		)
	}
}

func verifyBody(t *testing.T, body io.ReadCloser) (page int) {
	var q util.GraphQLRequest
	bb, err := io.ReadAll(body)
	if err != nil {
		t.Error("could not read request body:", err)
		return
	}
	if json.Unmarshal(bb, &q) != nil {
		t.Error("invalid GraphQL request body:", string(bb))
		return
	}
	om := util.MatchGQLParam("repository", "owner", q.Query)
	nm := util.MatchGQLParam("repository", "name", q.Query)
	cm := util.MatchGQLParam("releases", "after", q.Query)
	if om == nil || nm == nil || cm == nil {
		t.Error("wrong GraphQL query:", q.Query)
		return
	}

	owner, name, cur, ok := om[2], nm[2], cm[2], false
	if om[1] != "" {
		owner, ok = q.Variables[om[1]].(string)
		if !ok {
			t.Error("wrong GraphQL variable type")
		}
	}
	if nm[1] != "" {
		name, ok = q.Variables[nm[1]].(string)
		if !ok {
			t.Error("wrong GraphQL variable type")
		}
	}
	if owner != user {
		t.Error("wrong repo owner:", owner)
	}
	if name != repo {
		t.Error("wrong repo name:", name)
	}

	if cm[1] != "" {
		cur, ok = q.Variables[cm[1]].(string)
		if !ok {
			return
		}
	}
	if !strings.HasPrefix(cur, cursor) {
		t.Error("invalid cursor:", cur)
		return
	}
	page, err = strconv.Atoi(cur[len(cursor):])
	if err != nil {
		t.Error("invalid cursor:", cur)
	}
	return
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
