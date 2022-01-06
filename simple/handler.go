// Package simple implements the PEP 503 Simple Repository API
package simple

import (
	"net/http"
	"regexp"

	"github.com/plato-systems/pypihub/util"
)

// BaseURLPath is the endpoint of this service.
const BaseURLPath = "/simple/"

var pathSpec = regexp.MustCompile("^([a-z0-9-]*)/?$")

type handler struct {
	clients util.APIClientFactory
}

// NewHandler constructs a new Handler for this service.
func NewHandler(cf util.APIClientFactory) http.Handler {
	return &handler{cf}
}
