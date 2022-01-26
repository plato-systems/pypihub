package simple

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/util"
	"github.com/shurcooL/githubv4"
)

const (
	user, pass = "octocat", "123"
	pkg, repo  = "octopack", "test-repo"
	cursor     = "cceedd"
)

var assets = []ghAsset{
	{"Id0", pkg + "-1.1.1.tar.gz"},
	{"Id1", pkg + "-1.2.0.tar.gz"},
	{"Id2", pkg + "-1.2.0-py3-none-any.whl"},
}

func TestRoot(t *testing.T) {
	req, rec := setup("")
	req.SetBasicAuth(user, pass)

	makeNoreach(t).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestFoundPkg(t *testing.T) {
	testFound(t, pkg)
}

func TestFoundRepo(t *testing.T) {
	testFound(t, repo)
}

func testFound(t *testing.T, pkg string) {
	util.LoadConfigFile("./testdata/octo.toml")
	req, rec := setup(pkg)
	req.SetBasicAuth(user, pass)

	wraph(mockClient{t: t, assets: assets}).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Error("wrong status code: ", rec.Code)
	}

	bodyBuf, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatal("could not read response body: ", err)
	}
	body := string(bodyBuf)
	for _, a := range assets {
		anchor := fmt.Sprintf(
			`<a href="%s">%s</a>`,
			asset.MakeURL(a.ID, a.Name), a.Name,
		)
		if !strings.Contains(body, anchor) {
			t.Error("missing link:", anchor)
		}
	}
}

func TestNotFound(t *testing.T) {
	req, rec := setup(pkg)
	req.SetBasicAuth(user, pass)

	h := wraph(mockClient{t: t, err: errors.New("not found")})
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code: ", rec.Code)
	}
}

func TestUnauth(t *testing.T) {
	req, rec := setup(pkg)

	makeNoreach(t).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Error("wrong status code: ", rec.Code)
	}
}

func setup(pkg string) (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(
		http.MethodGet, path.Join(pathBase, pkg)+"/", nil,
	), httptest.NewRecorder()
}

func makeNoreach(t *testing.T) *handler {
	return wraph(mockClient{t: t, noreach: true})
}

func wraph(mc mockClient) *handler {
	return &handler{func(context.Context, string) util.GHv4Client {
		return mc
	}}
}

type mockClient struct {
	noreach bool
	assets  []ghAsset
	err     error
	t       *testing.T
}

func (m mockClient) Query(
	ctx context.Context,
	q interface{}, v map[string]interface{},
) error {
	if m.noreach {
		m.t.Error("should not call API")
		return m.err
	}

	query, ok := q.(*queryRepo)
	if !ok {
		m.t.Error("incorrect query type")
		return m.err
	}
	if len(m.assets) == 0 {
		return m.err
	}

	c, ok := v["cursor"].(*githubv4.String)
	i, err := 0, error(nil)
	if !ok {
		m.t.Error("incorrect cursor type")
		return m.err
	}
	if c != nil {
		i, err = strconv.Atoi(string(*c)[len(cursor):])
		if err != nil {
			m.t.Error("invalid cursor:", *c)
			return m.err
		}
	}

	rel := []ghRelease{{}}
	rel[0].ReleaseAssets.Nodes = m.assets[i : i+1]
	query.Repository.Releases.Nodes = rel

	pi := &query.Repository.Releases.PageInfo
	pi.HasNextPage = i+1 < len(m.assets)
	pi.EndCursor = fmt.Sprintf("%s%d", cursor, i+1)
	return m.err
}
