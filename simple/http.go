package simple

import (
	"log"
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/config"
)

var simpleRe = regexp.MustCompile("^/simple/([a-z0-9-]*)/?$")

func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	m := simpleRe.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	urlpath, pkg := m[0], m[1]
	if urlpath[len(urlpath)-1] != '/' {
		http.Redirect(w, r, urlpath+"/", http.StatusMovedPermanently)
		return
	}

	owner, token, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "401 unathorized", http.StatusUnauthorized)
		return
	}

	if pkg == "" {
		if err := tmplRoot.Execute(w, nil); err != nil {
			log.Println(err)
			error500(w)
		}
		return
	}

	repo := pkg
	for _, rep := range config.Conf.Replace {
		repo = rep.Re.ReplaceAllString(repo, rep.Repl)
	}
	assets, err := getRepoAssets(token, owner, repo)
	if err != nil {
		log.Println(err)
		error500(w)
		return
	}

	err = tmplPkg.Execute(w, argsTmplPkg{pkg, assets})
	if err != nil {
		log.Println(err)
		error500(w)
		return
	}
}

func error500(w http.ResponseWriter) {
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}
