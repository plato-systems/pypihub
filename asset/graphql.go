package asset

import (
	"context"

	"github.com/plato-systems/pypihub/util"
)

type queryAsset struct {
	Node struct {
		ReleaseAsset struct {
			URL string
		} `graphql:"... on ReleaseAsset"`
	} `graphql:"node(id: $assetID)"`
}

func getAssetURL(ctx context.Context, token, id string) (string, error) {
	client := util.NewGitHubv4Client(ctx, token)

	q, v := queryAsset{}, map[string]interface{}{
		"assetID": id,
	}
	err := client.Query(ctx, &q, v)
	if err != nil {
		return "", err
	}

	return q.Node.ReleaseAsset.URL, nil
}
