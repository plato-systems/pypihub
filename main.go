package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/simple"
	"github.com/plato-systems/pypihub/util"
)

var configFile = flag.String("c", "", "PyPIHub config file")

func main() {
	flag.Parse()
	if *configFile != "" {
		err := util.LoadConfigFile(*configFile)
		if err != nil {
			log.Fatal("[fatal] load config: ", err)
		}
	}

	asset.HandleHTTP()
	http.HandleFunc(simple.BaseURLPath, simple.ServeHTTP)

	s := util.Config.Server
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	log.Println("[info] Welcome to PyPIHub! Starting on", addr)

	if s.TLS.Crt == "" || s.TLS.Key == "" {
		log.Fatal(http.ListenAndServe(addr, nil))
	}
	log.Fatal(http.ListenAndServeTLS(addr, s.TLS.Crt, s.TLS.Key, nil))
}
