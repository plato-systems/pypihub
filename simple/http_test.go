package simple

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/plato-systems/pypihub/asset"
	"github.com/plato-systems/pypihub/util"
)

const (
	user, pass = "octocat", "123"
	pkg, repo  = "octopack", "test-repo"
	cursor     = "ccuu"
)

var assets = []ghAsset{
	{"Id0", pkg + "-1.1.1.tar.gz"},
	{"Id1", pkg + "-1.2.0.tar.gz"},
	{"Id2", pkg + "-1.2.0-py3-none-any.whl"},
}

func TestRoot(t *testing.T) {
	req, rec := setup("")
	req.SetBasicAuth(user, pass)

	wraph(makeUnreachableAPI(t)).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Error("wrong status code:", rec.Code)
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

	wraph(makeFoundAPI(t)).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Error("wrong status code:", rec.Code)
	}

	bodyBuf, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatal("could not read response body:", err)
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

	wraph(notFoundAPI).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Error("wrong status code:", rec.Code)
	}
}

func TestUnauth(t *testing.T) {
	req, rec := setup(pkg)

	wraph(makeUnreachableAPI(t)).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Error("wrong status code:", rec.Code)
	}
}
