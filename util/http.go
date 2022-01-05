package util

import "net/http"

// authOwners is the set of Basic Auth users allowed to use this server.
var authOwners map[string]bool

// ErrorHTTP replies to a request with an error code and
// the associated standard error text.
func ErrorHTTP(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

// AuthOwner checks a request for Basic Auth credentials
// that both exist and are allowed to use this server.
func AuthOwner(r *http.Request) (owner, token string, ok bool) {
	owner, token, ok = r.BasicAuth()
	if !ok || (len(authOwners) > 0 && !authOwners[owner]) {
		ok = false
		return
	}
	return
}
