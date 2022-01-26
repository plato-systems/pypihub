// Package simple implements the PEP 503 Simple Repository API.
package simple

import (
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/util"
)

const pathBase = "/simple/"

var pathSpec = regexp.MustCompile("^([a-z0-9-]*)/?$")

type handler struct {
	api util.APIClient
}

// HandleHTTP registers the Simple Repository API service in http.DefaultServeMux.
func HandleHTTP() {
	http.Handle(pathBase, &handler{util.GHv4Client{}})
}
