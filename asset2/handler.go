package asset2

import (
	"context"
	"net/http"
	"regexp"
)

// BaseURLPath is the endpoint of this service.
const BaseURLPath = "/asset/"

var pathSpec = regexp.MustCompile(`^([\w=]+)/[\w\.+-]+$`)

type handler struct {
	api API
}

func NewHandler(api API) http.Handler {
	if api == nil {
		return &handler{ghAPI{}}
	}
	return &handler{api}
}

type API interface {
	GetAsset(ctx context.Context, token, id string) (Asset, error)
}

type Asset struct {
	URL   string
	Owner string
}
