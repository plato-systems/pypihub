// Package asset implements PEP-503-compliant Release Assets redirects.
package asset

import (
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/util"
)

const pathBase = "/asset/"

var pathSpec = regexp.MustCompile(`^([\w=_\-]+)/[\w\.+\-]+$`)

type handler struct {
	makeGHv4Client util.GHv4ClientMaker
}

// HandleHTTP registers the Asset redirect service in http.DefaultServeMux.
func HandleHTTP() {
	http.Handle(pathBase, &handler{util.NewGHv4Client})
}
