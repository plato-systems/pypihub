package simple

type ghAsset struct {
	Name string
	URL  string
}

type ghRelease struct {
	Name          string
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
