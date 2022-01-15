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

func HandleHTTP() {
	http.Handle(pathBase, &handler{util.GHv4Client{}})
}
