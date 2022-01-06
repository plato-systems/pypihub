package asset

import (
	"log"
	"net/http"

	"github.com/plato-systems/pypihub/util"
)

// ServeHTTP redirects PEP-503-compliant file URLs to their temporary
// download URLs.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	client := h.clients.New(r.Context(), token)
	asset, err := getAsset(r.Context(), client, m[1])
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
