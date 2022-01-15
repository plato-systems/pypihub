package asset

import (
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/util"
)

const pathBase = "/asset/"

var pathSpec = regexp.MustCompile(`^([\w=]+)/[\w\.+-]+$`)

type handler struct {
	makeGHClient util.GHv4ClientMaker
}

func HandleHTTP() {
	http.Handle(pathBase, &handler{util.NewGHv4Client})
}
