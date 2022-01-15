// Package asset implements PEP-503-compliant Release Assets redirects.
package asset

import (
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/util"
)

const pathBase = "/asset/"

var pathSpec = regexp.MustCompile(`^([\w=]+)/[\w\.+-]+$`)

type handler struct {
	api util.APIClient
}

func HandleHTTP() {
	http.Handle(pathBase, &handler{util.GHv4Client{}})
}
