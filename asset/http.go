// Package asset implements PEP-503-compliant Release Assets redirects
package asset

import (
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/plato-systems/pypihub/util"
)

// BaseURLPath is the endpoint of this service.
const BaseURLPath = "/asset/"

var pathSpec = regexp.MustCompile(`^([\w=]+)/[\w\.+-]+$`)

// ServeHTTP redirects PEP-503-compliant file URLs to their temporary
// download URLs.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[info]", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		util.ErrorHTTP(w, http.StatusNotImplemented)
		return
	}

	m := pathSpec.FindStringSubmatch(r.URL.Path[len(BaseURLPath):])
	if m == nil {
		http.NotFound(w, r)
		return
	}

	owner, token, ok := util.AuthOwner(r)
	if !ok {
		util.ErrorHTTP(w, http.StatusUnauthorized)
		return
	}

	asset, err := getAsset(r.Context(), token, m[1])
	if err != nil {
		log.Printf("[warn] getAsset(%s): %v", m[0], err)
		http.NotFound(w, r)
		return
	}
	if asset.Release.Repository.Owner.Login != owner {
		// TODO: logging for suspicious activity?
		util.ErrorHTTP(w, http.StatusForbidden)
		return
	}

	http.Redirect(w, r, asset.URL, http.StatusTemporaryRedirect)
}

// MakeURL returns a redirect URL for the file with given ID and filename.
func MakeURL(id, filename string) string {
	return path.Join(BaseURLPath, id, filename)
}
