package vuetify

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// Generate returns generated vuetify components
func Generate(work util.GenerationWork, recipe *util.Recipe) {
	if !recipe.Vuetify.Generate {
		work.Waitgroup.Done()
		return
	}

	path := util.WorkingDir + "/web/" + recipe.Vuetify.App + "/src/bread"

	work.Waitgroup.Add(len(recipe.Entities) * 1) //2 jobs to be waited upon for each thread - Editor and List
	for _, entity := range recipe.Entities {
		if entity.Vuetify.NoGenerate {
			continue
		}

		go func(entity util.Entity) {
			var (
				data struct {
					Entity util.Entity
				}
			)

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

	menuEntities := func() []util.Entity {
		var items []util.Entity

		for i := range recipe.Entities {
			if !recipe.Entities[i].Vuetify.NotInMenu {
				items = append(items, recipe.Entities[i])
			}
		}

		return items
	}()

	output.GenerateAndSave(
		"Vuetify",
		"vuetify/js/routes.js.tmpl",
		path+"/routes.js",
		struct {
			Entities []util.Entity
		}{menuEntities},
		false,
		false,
	)

	// output.GenerateAndSave("Vuetify", "vuetify/store/index.js.tmpl", path+"/store/index.js", nil, true, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/actions.js.tmpl", path+"/store/actions.js", nil, true, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/getters.js.tmpl", path+"/store/getters.js", nil, true, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/mutations.js.tmpl", path+"/store/mutations.js", nil, true, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/types.js.tmpl", path+"/store/types.js", nil, true, false)

	work.Waitgroup.Done()
}
