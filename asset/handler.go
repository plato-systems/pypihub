// Package asset implements PEP-503-compliant Release Assets redirects
package asset

import (
	"net/http"
	"path"
	"regexp"

	"github.com/plato-systems/pypihub/util"
)

// BaseURLPath is the endpoint of this service.
const BaseURLPath = "/asset/"

var pathSpec = regexp.MustCompile(`^([\w=]+)/[\w\.+-]+$`)

type handler struct {
	clients util.APIClientFactory
}

// NewHandler constructs a new Handler for this service.
func NewHandler(cf util.APIClientFactory) http.Handler {
	return &handler{cf}
}

// MakeURL returns a redirect URL for the file with given ID and filename.
func MakeURL(id, filename string) string {
	return path.Join(BaseURLPath, id, filename)
}
