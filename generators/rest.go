package generators

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// GenerateREST returns generated code to run an http server
func GenerateREST(work util.GenerationWork, opts util.RestOpts, entities []util.Entity) error {
	work.Waitgroup.Add(len(entities) * 2) //2 jobs to be waited upon for each thread - _rest.go and _rest_hooks.go generation

	for _, entity := range entities {
		go func(entity util.Entity) {
			var (
				data struct {
					Package  string
					Entity   util.Entity
					Endpoint string
				}
			)

			if entity.Rest == nil {
				entity.Rest = &opts
			}

			data.Entity = entity
			data.Package = strings.ToLower(entity.Name)

			if opts.Prefix == "" {
				data.Endpoint = inflection.Plural(data.Package)
			} else {
				data.Endpoint = opts.Prefix + "/" + inflection.Plural(data.Package)
			}

			code, err := util.ExecuteTemplate("rest.go.tmpl", data)
			if entity.Rest.Hooks.PreCreate || entity.Rest.Hooks.PostCreate || entity.Rest.Hooks.PreRead || entity.Rest.Hooks.PostRead || entity.Rest.Hooks.PreList || entity.Rest.Hooks.PostList || entity.Rest.Hooks.PreUpdate || entity.Rest.Hooks.PostUpdate || entity.Rest.Hooks.PreDelete || entity.Rest.Hooks.PostDelete {
				hooks, e := util.ExecuteTemplate("rest_hooks.go.tmpl", struct {
					Hooks   util.RestHooks
					Name    string
					Package string
				}{entity.Rest.Hooks, entity.Name, data.Package})

				if e == nil {
					work.Done <- util.GeneratedCode{Generator: "GenerateRESTHooks", Code: hooks, Filename: fmt.Sprintf("models/%s/%s_rest_hooks.go", data.Package, data.Package), NoOverwrite: true}
				} else {
					work.Done <- util.GeneratedCode{Generator: "GenerateRESTHooks", Error: e}
				}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateRESTHooks", Error: util.ErrorSkip}
			}

			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateREST", Code: code, Filename: fmt.Sprintf("models/%s/%s_rest.go", data.Package, data.Package)}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateREST", Error: fmt.Errorf("failed to load execute template: %s", err)}
			}
		}(entity)
	}

	work.Waitgroup.Done()
	return nil
}