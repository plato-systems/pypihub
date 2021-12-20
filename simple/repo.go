package simple

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func getRepoAssets(token, owner, repo string) ([]ghAsset, error) {
	hc := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	client := githubv4.NewClient(hc)

	var assets []ghAsset
	q, v := queryRepo{}, map[string]interface{}{
		"repoOwner": githubv4.String(owner),
		"repoName":  githubv4.String(repo),
		"cursor":    (*githubv4.String)(nil),
	}
	for {
		err := client.Query(context.TODO(), &q, v)
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
