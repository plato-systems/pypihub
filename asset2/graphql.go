package asset2

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

func (g ghAPI) GetAsset(ctx context.Context, token, id string) (Asset, error) {
	client := util.NewGitHubv4Client(ctx, token)
	q, v := queryAsset{}, map[string]interface{}{"assetID": id}
	a, err := q.Node.ReleaseAsset, client.Query(ctx, &q, v)
	return Asset{URL: a.URL, Owner: a.Release.Repository.Owner.Login}, err
}
