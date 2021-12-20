package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/config"
	"github.com/plato-systems/pypihub/simple"
)

var configFile = flag.String("c", "", "PyPIHub config file")

func main() {
	flag.Parse()
	if *configFile != "" {
		err := config.LoadFile(*configFile)
		if err != nil {
			log.Fatal("[fatal] load config: ", err)
		}
	}

	http.HandleFunc(asset.BaseURLPath, asset.ServeHTTP)
	http.HandleFunc(simple.BaseURLPath, simple.ServeHTTP)

	addr := fmt.Sprintf("%s:%d", config.Conf.Host, config.Conf.Port)
	log.Println("[info] Welcome to PyPIHub! Starting on", addr)

	tls := config.Conf.TLS
	if tls.Cert == "" || tls.Key == "" {
		log.Fatal(http.ListenAndServe(addr, nil))
	}
	log.Fatal(http.ListenAndServeTLS(addr, tls.Cert, tls.Key, nil))
}
