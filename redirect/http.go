// Package redirect implements a multi-purpose caching redirect layer
package redirect

import (
	"log"
	"net/http"
)

const BaseURLPath = "/redirect/"

func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[info]", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "501 not implemented", http.StatusNotImplemented)
		return
	}

	user, pass, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "401 unauthorized", http.StatusUnauthorized)
		return
	}

	e := table[r.URL.Path[len(BaseURLPath):]]
	if e == nil {
		http.NotFound(w, r)
		return
	}
	if user != e.user || pass != e.pass {
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	http.Redirect(w, r, e.dest, http.StatusTemporaryRedirect)
}
