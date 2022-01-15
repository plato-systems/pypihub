package simple

import (
	"context"
	"net/http"
	"regexp"
)

const pathBase = "/simple/"

var pathSpec = regexp.MustCompile("^([a-z0-9-]*)/?$")

type handler struct {
	api interface {
		getRepoAssets(ctx context.Context, token, owner, repo string) ([]ghAsset, error)
	}
}

type ghAPI struct{}

func HandleHTTP() {
	http.Handle(pathBase, &handler{ghAPI{}})
}
