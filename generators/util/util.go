package util

import (
	"github.com/fluxynet/gocipe/output"
	utils "github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(out output.Output) {
	data := struct{ AppImportPath string }{utils.AppImportPath}

	out.GenerateAndOverwrite("Util Rice", "util/rice.go.tmpl", "util/rice.gocipe.go", nil)
	out.GenerateAndOverwrite("Util Credentials", "util/credentials.go.tmpl", "util/credentials/credentials.gocipe.go", nil)
	out.GenerateAndOverwrite("Util", "util/util.go.tmpl", "util/util.gocipe.go", nil)
	out.GenerateAndOverwrite("Util Web", "util/web.go.tmpl", "util/web/web.gocipe.go", nil)
	out.GenerateAndOverwrite("Util Ws", "util/ws.go.tmpl", "util/web/ws.gocipe.go", nil)
	out.GenerateAndOverwrite("Util Fileupload", "util/files.go.tmpl", "util/files/files.gocipe.go", data)
	out.GenerateAndOverwrite("Util Imagist", "util/imagist.go.tmpl", "util/imagist/imagist.gocipe.go", data)
}
