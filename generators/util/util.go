package util

import (
	"github.com/fluxynet/gocipe/output"

	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(work util.GenerationWork) {
	data := struct{ AppImportPath string }{util.AppImportPath}

	output.GenerateAndSave("Util", "util/rice.go.tmpl", "util/rice.gocipe.go", nil, false)
	output.GenerateAndSave("Util", "util/credentials.go.tmpl", "util/credentials/credentials.gocipe.go", nil, false)
	output.GenerateAndSave("Util", "util/util.go.tmpl", "util/util.gocipe.go", nil, false)
	output.GenerateAndSave("Util", "util/web.go.tmpl", "util/web/web.gocipe.go", nil, false)
	output.GenerateAndSave("Util", "util/ws.go.tmpl", "util/web/ws.gocipe.go", nil, false)
	output.GenerateAndSave("Util", "util/files.go.tmpl", "util/files/files.gocipe.go", data, false)
	output.GenerateAndSave("Util", "util/imagist.go.tmpl", "util/imagist/imagist.gocipe.go", data, false)
	output.GenerateAndSave("Util", "util/grpc.go.tmpl", "util/grpcx/grpc.gocipe.go", nil, false)
	work.Waitgroup.Done()
}
