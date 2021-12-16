package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/config"
)

var configFile = flag.String("c", "", "PyPIHub config file")

func main() {
	flag.Parse()
	if *configFile != "" {
		err := config.LoadFile(*configFile)
		if err != nil {
			log.Fatal("failed to load config: ", err)
		}
	}

	http.HandleFunc("/simple/", handleSimple)

	addr := fmt.Sprintf("%s:%d", config.Conf.Host, config.Conf.Port)
	tls := config.Conf.TLS
	if tls.Cert == "" || tls.Key == "" {
		log.Fatal(http.ListenAndServe(addr, nil))
	}
	log.Fatal(http.ListenAndServeTLS(addr, tls.Cert, tls.Key, nil))
}

var simpleRe = regexp.MustCompile("^/simple/([a-z0-9-]*)/?$")

func handleSimple(w http.ResponseWriter, r *http.Request) {
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

	user, token, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "401 unathorized", http.StatusUnauthorized)
		return
	}

	if pkg == "" {
		fmt.Fprintln(w, "<h1>Welcome to PyPIHub!</h1>")
		return
	}

	fmt.Fprintln(w, "list pkg:", pkg, user, token)
}
