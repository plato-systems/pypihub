package asset

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type queryAsset struct {
	Node struct {
		ReleaseAsset struct {
			URL string
		} `graphql:"... on ReleaseAsset"`
	} `graphql:"node(id: $assetID)"`
}

func getAssetURL(id, token string) (string, error) {
	hc := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	client := githubv4.NewClient(hc)

	q, v := queryAsset{}, map[string]interface{}{
		"assetID": id,
	}
	err := client.Query(context.TODO(), &q, v)
	if err != nil {
		return "", err
	}

	return q.Node.ReleaseAsset.URL, nil
}
