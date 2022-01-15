package simple

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/util"
)

const (
	user, pass = "octocat", "123"
	pkg, repo  = "octopack", "test-repo"
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

	h := handler{mockAPI{assets: assets}}
	h.ServeHTTP(rec, req)
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

	h := handler{mockAPI{err: errors.New("not found")}}
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
	return &handler{mockAPI{noreach: t}}
}

type mockAPI struct {
	noreach *testing.T
	assets  []ghAsset
	err     error
}

func (m mockAPI) getRepoAssets(ctx context.Context, token, owner, repo string) ([]ghAsset, error) {
	if m.noreach != nil {
		m.noreach.Error("should not call API")
	}
	return m.assets, m.err
}
