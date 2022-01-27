package asset

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plato-systems/pypihub/util"
)

func setup() (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(
		http.MethodGet, MakeURL(id, file), nil,
	), httptest.NewRecorder()
}

func wraph(serve http.HandlerFunc) http.Handler {
	return &handler{util.NewGHv4ClientMaker(serve)}
}

func makeFoundAPI(t *testing.T) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		verifyBody(t, r.Body)
		fmt.Fprintf(rw, `{"data": { "node": {
			"url": "%s",
			"release": { "repository": { "owner": {
				"login": "%s"
			}}}
		}}}`, location, user)
	}
}

func verifyBody(t *testing.T, body io.ReadCloser) {
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

	m := util.MatchGQLParam("node", "id", q.Query)
	if m == nil {
		t.Error("wrong GraphQL query:", q.Query)
		return
	}

	assetID, ok := m[2], false
	if m[1] != "" { // variable
		assetID, ok = q.Variables[m[1]].(string)
		if !ok {
			t.Error("wrong GraphQL variable type")
			return
		}
	}
	if assetID != id {
		t.Error("wrong asset id:", assetID)
	}
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
