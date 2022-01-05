// Package simple implements the PEP 503 Simple Repository API
package simple

import (
	"log"
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/util"
)

// BaseURLPath is the endpoint of this service.
const BaseURLPath = "/simple/"

var pathSpec = regexp.MustCompile("^([a-z0-9-]*)/?$")

// ServeHTTP lists downloadable files for requested Package
// from Release Assets of its hosting Repo.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Println("[info]", r.Method, path)
	if r.Method != http.MethodGet {
		util.ErrorHTTP(w, http.StatusNotImplemented)
		return
	}

	m := pathSpec.FindStringSubmatch(path[len(BaseURLPath):])
	if m == nil {
		http.NotFound(w, r)
		return
	}
	if path[len(path)-1] != '/' {
		http.Redirect(w, r, path+"/", http.StatusMovedPermanently)
		return
	}

	owner, token, ok := util.AuthOwner(r)
	if !ok {
		util.ErrorHTTP(w, http.StatusUnauthorized)
		return
	}

	pkg := m[1]
	if pkg == "" { // GET /simple/
		w.Write([]byte(htmlRoot))
		return
	}

	repo := pkg
	for _, rep := range util.Config.GitHub.Replace {
		repo = rep.Patt.ReplaceAllString(repo, rep.Repl)
	}

	assets, err := getRepoAssets(r.Context(), token, owner, repo)
	if err != nil {
		log.Printf("[warn] getRepoAssets(%s/%s): %v", owner, repo, err)
		http.NotFound(w, r)
		return
	}

	err = tmplPkg.Execute(w, argsPkg{pkg, assets})
	if err != nil {
		log.Printf("[error] tmplPkg.Execute(%s): %v", pkg, err)
		util.ErrorHTTP(w, http.StatusInternalServerError)
		return
	}
}
