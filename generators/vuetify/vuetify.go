package vuetify

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// Generate returns generated vuetify components
func Generate(work util.GenerationWork, opts util.VuetifyOpts, entities []util.Entity) {
	if !opts.Generate {
		work.Waitgroup.Done()
		return
	}

	path := util.WorkingDir + "/web/" + opts.App + "/src/bread"

	work.Waitgroup.Add(len(entities) * 1) //2 jobs to be waited upon for each thread - Editor and List
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

			filename := path + "/views/" + inflection.Plural(data.Entity.Name)
			list, err := util.ExecuteTemplate("vuetify/views/list.vue.tmpl", data)

			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyList", Code: list, Filename: filename + "List.vue"}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyList", Error: err, GeneratedHeaderFormat: "<-- %s -->"}
			}

			// edit, err := util.ExecuteTemplate("vuetify_edit.vue.tmpl", data)
			// if err == nil {
			// 	work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyEdit", Code: edit, Filename: filename + "Edit.vue"}
			// } else {
			// 	work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyEdit", Error: err, GeneratedHeaderFormat: "<-- %s -->"}
			// }
		}(entity)
	}

	output.GenerateAndSave("Vuetify", "vuetify/store/routes.js.tmpl", path+"/store/routes.js", nil, true, false)
	output.GenerateAndSave("Vuetify", "vuetify/store/index.js.tmpl", path+"/store/index.js", nil, true, false)
	output.GenerateAndSave("Vuetify", "vuetify/store/actions.js.tmpl", path+"/store/actions.js", nil, true, false)
	output.GenerateAndSave("Vuetify", "vuetify/store/getters.js.tmpl", path+"/store/getters.js", nil, true, false)
	output.GenerateAndSave("Vuetify", "vuetify/store/mutations.js.tmpl", path+"/store/mutations.js", nil, true, false)
	output.GenerateAndSave("Vuetify", "vuetify/store/types.js.tmpl", path+"/store/types.js", nil, true, false)

	work.Waitgroup.Done()
}
