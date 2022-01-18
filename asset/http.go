// Package asset implements PEP-503-compliant Release Assets redirects.
package asset

import (
	"log"
	"net/http"
	"path"

	"github.com/plato-systems/pypihub/util"
)

// ServeHTTP redirects artificial file URLs to their temporary download URLs.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[info]", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		util.ErrorHTTP(w, http.StatusNotImplemented)
		return
	}

	m := pathSpec.FindStringSubmatch(r.URL.Path[len(pathBase):])
	if m == nil {
		http.NotFound(w, r)
		return
	}

	owner, token, ok := util.AuthOwner(r)
	if !ok {
		util.ErrorHTTP(w, http.StatusUnauthorized)
		return
	}

	client := h.makeGHv4Client(r.Context(), token)
	a, err := getAsset(r.Context(), client, m[1])
	if err != nil {
		log.Printf("[warn] getAsset(%s): %v", m[0], err)
		http.NotFound(w, r)
		return
	}
	if a.Release.Repository.Owner.Login != owner {
		// TODO: logging for suspicious activity?
		util.ErrorHTTP(w, http.StatusForbidden)
		return
	}

	http.Redirect(w, r, a.URL, http.StatusFound)
}

// MakeURL produces an artificial file URL for an Asset.
func MakeURL(id, filename string) string {
	return path.Join(pathBase, id, filename)
}
