// Package redirect implements a multi-purpose caching redirect layer
package redirect

import (
	"net/http"
)

const BaseURLPath = "/redirect/"

func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	e := table[r.URL.Path[len(BaseURLPath):]]
	if e == nil {
		http.NotFound(w, r)
		return
	}

	user, pass, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "401 unauthorized", http.StatusUnauthorized)
		return
	}
	if user != e.user || pass != e.pass {
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	http.Redirect(w, r, e.dest, http.StatusTemporaryRedirect)
}
