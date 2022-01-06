package util_test

import (
	"net/http"
	"testing"

	"github.com/plato-systems/pypihub/util"
)

const (
	user = "octocat"
	pass = "test123"
)

func newRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth(user, pass)
	return req
}

func TestUnauth(t *testing.T) {
	util.LoadConfigFile("./testdata/none.toml")
	r := http.Request{}
	if _, _, ok := util.AuthOwner(&r); ok {
		t.Error("unauthenticated request accepted")
	}
}

func TestPublic(t *testing.T) {
	util.LoadConfigFile("./testdata/none.toml")
	owner, token, ok := util.AuthOwner(newRequest())
	if !ok {
		t.Error("authenticated request not accepted")
	}
	if owner != user {
		t.Error("incorrect owner:", owner)
	}
	if token != pass {
		t.Error("incorrect token:", token)
	}
}

func TestOne(t *testing.T) {
	util.LoadConfigFile("./testdata/one.toml")

	owner, token, ok := util.AuthOwner(newRequest())
	if !ok {
		t.Error("correctly authenticated request not accepted")
	}
	if owner != user {
		t.Error("incorrect owner:", owner)
	}
	if token != pass {
		t.Error("incorrect token:", token)
	}
}

func TestBadOne(t *testing.T) {
	util.LoadConfigFile("./testdata/one.toml")
	r, user := newRequest(), "nobody"
	r.SetBasicAuth(user, pass)

	owner, token, ok := util.AuthOwner(r)
	if ok {
		t.Error("incorrectly authenticated request accepted")
	}
	if owner != user {
		t.Error("incorrect owner:", owner)
	}
	if token != pass {
		t.Error("incorrect token:", token)
	}
}

func TestTwo(t *testing.T) {
	util.LoadConfigFile("./testdata/two.toml")
	r := newRequest()

	owner, token, ok := util.AuthOwner(r)
	if !ok {
		t.Error("correctly authenticated request not accepted")
	}
	if owner != user {
		t.Error("incorrect owner:", owner)
	}
	if token != pass {
		t.Error("incorrect token:", token)
	}

	user := "octorg"
	r.SetBasicAuth(user, pass)
	owner, token, ok = util.AuthOwner(r)
	if !ok {
		t.Error("correctly authenticated request not accepted")
	}
	if owner != user {
		t.Error("incorrect owner:", owner)
	}
	if token != pass {
		t.Error("incorrect token:", token)
	}
}

func TestBadTwo(t *testing.T) {
	util.LoadConfigFile("./testdata/two.toml")
	r, user := newRequest(), "nobody"
	r.SetBasicAuth(user, pass)

	owner, token, ok := util.AuthOwner(r)
	if ok {
		t.Error("incorrectly authenticated request accepted")
	}
	if owner != user {
		t.Error("incorrect owner:", owner)
	}
	if token != pass {
		t.Error("incorrect token:", token)
	}
}
