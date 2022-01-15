package asset

import (
	"context"

	"github.com/plato-systems/pypihub/util"
)

type queryAsset struct {
	Node struct {
		ReleaseAsset ghAsset `graphql:"... on ReleaseAsset"`
	} `graphql:"node(id: $assetID)"`
}

type ghAsset struct {
	URL     string
	Release struct {
		Repository struct {
			Owner struct {
				Login string
			}
		}
	}
}

type ghAPI struct{}

// TODO: return more meaningful errors
func (g ghAPI) getAsset(ctx context.Context, token, id string) (asset, error) {
	client := util.NewGitHubv4Client(ctx, token)
	q, v := queryAsset{}, map[string]interface{}{"assetID": id}
	a, err := q.Node.ReleaseAsset, client.Query(ctx, &q, v)
	return asset{a.URL, a.Release.Repository.Owner.Login}, err
}
