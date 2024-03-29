package simple

import (
	"context"

	"github.com/shurcooL/githubv4"
)

type ghAsset struct {
	ID   string
	Name string
}

type ghRelease struct {
	ReleaseAssets struct {
		Nodes []ghAsset
	} `graphql:"releaseAssets(first: 32)"`
}

type queryRepo struct {
	Repository struct {
		Releases struct {
			Nodes    []ghRelease
			PageInfo struct {
				EndCursor   string
				HasNextPage bool
			}
		} `graphql:"releases(first: 64, after: $cursor)"`
	} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
}

// TODO: return more meaningful errors
func getRepoAssets(
	ctx context.Context,
	client *githubv4.Client,
	owner, repo string,
) ([]ghAsset, error) {
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
