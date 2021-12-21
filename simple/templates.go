package simple

import (
	"html/template"

	"github.com/plato-systems/pypihub/asset"
)

const htmlRoot = `<!DOCTYPE html>
<html>
<head>
	<title>PyPIHub: root</title>
</head>
<body>
	<h1>Welcome to PyPIHub!</h1>
</body>
</html>
`

var tmplPkg = template.Must(template.New("pkg").Funcs(template.FuncMap{
	"assetURL": asset.MakeURL,
}).Parse(`<!DOCTYPE html>
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
`))

type argsPkg struct {
	Name   string
	Assets []ghAsset
}
