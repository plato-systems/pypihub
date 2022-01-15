package asset

import (
	"context"
	"net/http"
	"regexp"
)

const pathBase = "/asset/"

var pathSpec = regexp.MustCompile(`^([\w=]+)/[\w\.+-]+$`)

type asset struct {
	url   string
	owner string
}

type handler struct {
	api interface {
		getAsset(ctx context.Context, token, id string) (asset, error)
	}
}

type ghAPI struct{}

func HandleHTTP() {
	http.Handle(pathBase, &handler{ghAPI{}})
}
