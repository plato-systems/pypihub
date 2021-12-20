package simple

import (
	"html/template"

	"github.com/plato-systems/pypihub/asset"
)

var (
	htmlRoot = `<!DOCTYPE html>
<html>
<head>
	<title>PyPIHub: root</title>
</head>
<body>
	<h1>Welcome to PyPIHub!</h1>
</body>
</html>
`

	htmlPkg = `<!DOCTYPE html>
<html>
<head>
	<title>PyPIHub: package {{.Name}}</title>
</head>
<body>
	<h1>Links for {{.Name}}</h1>
{{range .Assets}}
	<a href="{{assetURL .ID .Name}}">{{.Name}}</a><br />
{{end}}
</body>
</html>
`
)

var (
	tmplRoot = template.Must(template.New("root").Parse(htmlRoot))
	tmplPkg  = template.Must(template.New("pkg").Funcs(template.FuncMap{
		"assetURL": asset.MakeURL,
	}).Parse(htmlPkg))
)

type argsTmplPkg struct {
	Name   string
	Assets []ghAsset
}
