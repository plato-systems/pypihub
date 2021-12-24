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
Expect Basic Auth with requested project-Repo owner's login as user
and invoking user's `repo`-scoped GitHub PAT as password

All GitHub API operations are authenticated using the PAT (Basic Auth
password) allowing for the access-control policies defined on GitHub users
to apply uniformly: a user can only access a Repo's Assets as Package files
if he has access to Repo itself (since his PAT is used through the API)

#### `GET /simple/<project>/`
Find and return links to all files for given `project`:
1. Convert param `project` to Repo name, set Repo owner to Basic Auth user
2. Collect Assets from all Releases for Repo (max 32 Assets per Release)
3. Emit `/asset/<ID>/<name>` anchor tags for each Asset

#### `GET /asset/<ID>/<name>`
Redirect to download file of Release Asset with given `ID`:
1. Look up Release Asset with param `ID` (globally unique ID)
2. Verify associated Repo belongs to Basic Auth user
3. Redirect to temporary download URL for Asset file

Required for PEP 503 compliance since GitHub API download URLs for Assets
do not include the filename as the final path component

### Non-features
* Package search: only project-to-Repo mapping supported, not vice-versa
  * Thus, the root URL (`GET /simple/`) returns no projects
  * Required for full PEP 503 compliance :(
* Package upload: *might* become supported in future versions
  * Not required for PEP 503 compliance ;)

## Installation
From source:
```
$ go install github.com/plato-systems/pypihub@latest
```

## Usage
View help:
```
$ pypihub -h
```

### Example
Suppose your GitHub account is `ocotocat` and you want to access some
Python Packages hosted in Repos (you have access to) belonging to `octorg`

First, [create a PAT](https://github.com/settings/tokens) with `repo` scope

Then, run a PyPIHub server with default config:
```
$ pypihub
```
TODO: more coming soon!

## Configuration
TODO: coming soon!
### Default
TODO: coming soon!

## Related
* [pywharf](https://github.com/pywharf/pywharf): Inspiration for PyPIHub
  * PyPI server using one GitHub Repo for all hosted projects
    rather than per-project Repos
  * Supports Package search through special index file
  * Supports Package upload (not confirmed)
