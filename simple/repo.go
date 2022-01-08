// Package simple implements the PEP 503 Simple Repository API
package simple

import (
	"context"

	"github.com/plato-systems/pypihub/util"
	"github.com/shurcooL/githubv4"
)

func getRepoAssets(
	ctx context.Context,
	token, owner, repo string,
) ([]ghAsset, error) {
	client := util.NewGitHubv4Client(ctx, token)
	assets := []ghAsset{}
	q, v := queryRepo{}, map[string]interface{}{
		"repoOwner": githubv4.String(owner),
		"repoName":  githubv4.String(repo),
		"cursor":    (*githubv4.String)(nil),
	}
	for {
		err := client.Query(ctx, &q, v)
		if err != nil {
			return nil, err
		}

		for _, rel := range q.Repository.Releases.Nodes {
			assets = append(assets, rel.ReleaseAssets.Nodes...)
		}

		if !q.Repository.Releases.PageInfo.HasNextPage {
			break
		}
		v["cursor"] = githubv4.String(q.Repository.Releases.PageInfo.EndCursor)
	}

	return assets, nil
}
