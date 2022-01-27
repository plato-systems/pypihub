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

	owner, oerr := util.MatchGQLParam("repository", "owner", q)
	name, nerr := util.MatchGQLParam("repository", "name", q)
	cur, cerr := util.MatchGQLParam("releases", "after", q)
	if oerr != nil || nerr != nil {
		t.Error("wrong GraphQL query:", q.Query)
		return
	}
	if owner != user {
		t.Error("wrong repo owner:", owner)
	}
	if name != repo {
		t.Error("wrong repo name:", name)
	}

	if cerr != nil {
		return
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
