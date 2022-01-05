# PyPIHub
Python Package Index backed by GitHub Releases

## Overview
An HTTP server (mostly) implementing the
[PEP 503](https://www.python.org/dev/peps/pep-0503/)
Simple Repository API and using the
[GitHub GraphQL API](https://docs.github.com/v4/)
to find Python Packages hosted as GitHub Release Assets

See an informative usage [example](#example) below!

### Architecture
![architecture diagram](doc/arch.png)

### Endpoints
Expect Basic Auth with:
* user: Owner of requested Repo
* password: `repo`-scoped GitHub PAT of invoking user

All GitHub API operations are authenticated using the PAT, allowing
for the access-control policies defined on GitHub Accounts to apply
uniformly: a user can only access a Repo's Assets as Package files if
he has access to Repo itself (since his PAT is used through the API)

#### `GET /simple/<pkg>/`
Find and return links to all files for named Python Package:
1. Convert param `pkg` to Repo name, set Repo owner to Basic Auth user
2. Collect Assets from all Releases for Repo (max 32 Assets per Release)
3. Emit `/asset/<ID>/<name>` anchor tags for each Asset

#### `GET /asset/<ID>/<name>`
Redirect to download file of Release Asset with given `ID`:
1. Look up Release Asset with param `ID` (globally unique)
2. Verify associated Repo belongs to Basic Auth user
3. Redirect to temporary download URL for Asset file

Required for PEP 503 compliance since GitHub API download URLs for Assets
do not include the filename as the final path component

### Non-features
* Package search: only Package-to-Repo mapping supported, not vice-versa
  * Thus, the root URL (`GET /simple/`) returns no Packages
  * Required for full PEP 503 compliance :(
* Package upload: *might* become supported in future versions
  * Not required for PEP 503 compliance ;)

## Installation
From source:
```sh
$ go install github.com/plato-systems/pypihub@latest
```

## Usage
View help:
```sh
$ pypihub -h
```

### Example
Suppose your GitHub Account is `octocat` and you want to access a
Python Package named `octopack` hosted in a Repo of the same name
(that you have access to) belonging to the `octorg` Organization

First, [create a PAT](https://github.com/settings/tokens) for your
`octocat` Account with `repo` scope

Then, run a PyPIHub server with default config:
```sh
$ pypihub
```

Now, suppose the `octopack` Repo has a Release with an Asset named
`octopack-1.2.3.tar.gz` indicating a source distribution of the
`octopack` Python Package with version `1.2.3`

To install it (from the same host running the PyPIHub server):
```sh
$ pip install -i http://octorg:${PAT}@localhost:3141/simple/ octopack==1.2.3
```
Notice the Basic Auth user is specified as `octorg` *not* `octocat`!
This is important since `octorg` owns the `octopack` Repo, not `octocat`.
The PAT however, *is* `octocat`'s since that is the user requesting the
Assets.

## Configuration
Specified in a [TOML](https://toml.io/en/) file optionally passed
to the `pypihub` binary

Here is an annotated example configuration
(not necessarily comprehensive; refer to `util/config.go`):
```toml
[server]
host = "127.0.0.1"  # listen only on loopback interface
port = 8080         # listen on port 8080

[server.tls]        # enable HTTPS
crt = "tls.crt"     # use TLS certificate in file tls.crt
key = "tls.key"     # use TLS private key in file tls.key

[github]
# Serve only Packages hosted in Repos belonging to these Owners
owners = [ "octorg", "octocat" ]
# Otherwise, anyone could spam your server with requests for Packages
# in public Repos for example, when all you want it to do is serve
# Packages in Repos owned by either 'octorg' or 'octocat'

# Transform Package names to Repo names through a series of
# string replacements applied to all Package names processed
[[github.replace]]  # strip 'octo-' prefix
patt = "^octo-"     # match pattern specified as Regular Expression
repl = ""           # Go-style replacement string

[[github.replace]]  # replace all dashes with underscores
patt = "-"
repl = "_"
# Replace directives are applied in the order given, e.g.:
# (Package) 'octo-hello-world' => 'hello-world' => (Repo) 'hello_world'
```
All keys are optional and when not given follow the [Defaults](#defaults)

### Defaults
* No HTTPS
* Listens on port `3141` on all interfaces
* All Repo Owners allowed (public instance)
* Package name is Repo name

## Related
* [pywharf](https://github.com/pywharf/pywharf): Inspiration for PyPIHub
  * PyPI server using one GitHub Repo for all hosted Packages
    rather than per-Package Repos
  * Supports Package search through special index file
  * Supports Package upload (not confirmed)
