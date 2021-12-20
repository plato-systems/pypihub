// Package asset implements PEP-503-compliant Release Assets redirects
package asset

import (
	"log"
	"net/http"
	"path"
	"regexp"
)

const BaseURLPath = "/asset/"

var pathSpec = regexp.MustCompile(`^([\w=]+)/[\w\.+-]+$`)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[info]", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "501 not implemented", http.StatusNotImplemented)
		return
	}

	m := pathSpec.FindStringSubmatch(r.URL.Path[len(BaseURLPath):])
	if m == nil {
		http.NotFound(w, r)
		return
	}

	_, token, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "401 unathorized", http.StatusUnauthorized)
		return
	}

	url, err := getAssetURL(m[1], token)
	if err != nil {
		log.Printf("[warn] getAssetURL(%s): %v", m[0], err)
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func MakeURL(id, filename string) string {
	return path.Join(BaseURLPath, id, filename)
}
