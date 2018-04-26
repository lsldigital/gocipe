package generators

import (
	"strings"

	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// GenerateVuetify returns generated vuetify components
func GenerateVuetify(work util.GenerationWork, restOpts util.RestOpts, opts util.VuetifyOpts, entities []util.Entity) error {
	if !opts.Generate {
		work.Waitgroup.Done()
		return nil
	}

	if opts.Module == "" {
		opts.Module = "web"
	}

	work.Waitgroup.Add(len(entities) * 2) //2 jobs to be waited upon for each thread - Editor and List
	for _, entity := range entities {
		go func(entity util.Entity) {
			var (
				data struct {
					Endpoint string
					Entity   util.Entity
					Prefix   string
				}
			)

			if entity.Vuetify == nil {
				entity.Vuetify = &opts
			}

			data.Entity = entity
			data.Endpoint = restOpts.Prefix + inflection.Plural(strings.ToLower(entity.Name))
			data.Prefix = restOpts.Prefix
			filename := opts.Module + "/src/modules/gocipe/views/" + inflection.Plural(data.Entity.Name)
			list, err := util.ExecuteTemplate("vuetify_list.vue.tmpl", data)

			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyList", Code: list, Filename: filename + "List.vue"}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyList", Error: err}
			}

			edit, err := util.ExecuteTemplate("vuetify_edit.vue.tmpl", data)
			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyEdit", Code: edit, Filename: filename + "Edit.vue"}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyEdit", Error: err}
			}
		}(entity)
	}
	work.Waitgroup.Add(6) //6 stubs

	var (
		stub string
		err  error
	)
	path := opts.Module + "/src/modules/gocipe/store/"

	stub, err = util.ExecuteTemplate("vuetify_routes.js.tmpl", struct {
		Entities []util.Entity
	}{entities})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleRoutes", Code: stub, Filename: path + "routes.js"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleRoutes", Error: err}
	}

	stub, err = util.ExecuteTemplate("vuetify_index.js.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleIndex", Code: stub, Filename: path + "index.js", NoOverwrite: true}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleIndex", Error: err}
	}

	stub, err = util.ExecuteTemplate("vuetify_actions.js.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleActions", Code: stub, Filename: path + "actions.js", NoOverwrite: true}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleActions", Error: err}
	}

	stub, err = util.ExecuteTemplate("vuetify_getters.js.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleGetters", Code: stub, Filename: path + "getters.js", NoOverwrite: true}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleGetters", Error: err}
	}

	stub, err = util.ExecuteTemplate("vuetify_mutations.js.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleMutations", Code: stub, Filename: path + "mutations.js", NoOverwrite: true}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleMutations", Error: err}
	}

	stub, err = util.ExecuteTemplate("vuetify_types.js.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleTypes", Code: stub, Filename: path + "types.js", NoOverwrite: true}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyModuleTypes", Error: err}
	}

	work.Waitgroup.Done()
	return nil
}
