package util

import "net/http"

var authOwners map[string]bool

func ErrorHTTP(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func AuthOwner(r *http.Request) (owner, token string, ok bool) {
	owner, token, ok = r.BasicAuth()
	if !ok || (len(authOwners) > 0 && !authOwners[owner]) {
		ok = false
		return
	}
	return
}
