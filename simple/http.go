// Package simple implements the PEP 503 Simple Repository API
package simple

import (
	"log"
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/config"
	"github.com/plato-systems/pypihub/redirect"
)

const BaseURLPath = "/simple/"

var pathSpec = regexp.MustCompile("^" + BaseURLPath + "([a-z0-9-]*)/?$")

func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	m := pathSpec.FindStringSubmatch(path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	if path[len(path)-1] != '/' {
		http.Redirect(w, r, path+"/", http.StatusMovedPermanently)
		return
	}

	owner, token, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "401 unathorized", http.StatusUnauthorized)
		return
	}

	pkg := m[1]
	if pkg == "" { // GET /simple/
		if err := tmplRoot.Execute(w, nil); err != nil {
			log.Println("tmplRoot.Execute(): ", err)
			error500(w)
		}
		return
	}

	repo := pkg
	for _, rep := range config.Conf.Replace {
		repo = rep.Re.ReplaceAllString(repo, rep.Repl)
	}

	assets, err := repoGetAssets(token, owner, repo)
	if err != nil {
		log.Printf("repoGetAssets(%s/%s): %v", owner, repo, err)
		http.NotFound(w, r)
		return
	}

	for i := range assets {
		a := &assets[i]
		a.URL = redirect.Register(a.Name, a.URL, owner, token)
	}

	err = tmplPkg.Execute(w, argsTmplPkg{pkg, assets})
	if err != nil {
		log.Printf("tmplPkg.Execute(%s): %v", pkg, err)
		error500(w)
		return
	}
}

func error500(w http.ResponseWriter) {
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}
