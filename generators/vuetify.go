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
			} else if entity.Vuetify.Path == "" {
				entity.Vuetify.Path = opts.Path
			}

			path := entity.Vuetify.Path
			if path == "" {
				path = "vuetify"
			}

			data.Entity = entity
			data.Endpoint = restOpts.Prefix + "/" + inflection.Plural(strings.ToLower(entity.Name))
			data.Prefix = restOpts.Prefix
			filename := path + "/" + data.Entity.Name
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

	work.Waitgroup.Done()
	return nil
}
