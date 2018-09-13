package vuetify

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// Generate returns generated vuetify components
func Generate(work util.GenerationWork, recipe *util.Recipe, entities map[string]util.Entity) {
	if !recipe.Vuetify.Generate {
		work.Waitgroup.Done()
		return
	}

	path := util.WorkingDir + "/web/" + recipe.Vuetify.App + "/src/gocipe"

	// work.Waitgroup.Add(len(recipe.Entities) * 1) //2 jobs to be waited upon for each thread - Editor and List
	for _, entity := range entities {
		if entity.Vuetify.NoGenerate {
			continue
		}

		// go func(entity util.Entity) {
		var (
			data struct {
				Entity   util.Entity
				Entities map[string]util.Entity
			}
		)

		data.Entity = entity
		data.Entities = entities

		filename := path + "/views/" + inflection.Plural(data.Entity.Name)

		output.GenerateAndSave(
			"VuetifyList",
			"vuetify/views/list.vue.tmpl",
			filename+"List.vue",
			data,
			false,
		)

		output.GenerateAndSave(
			"VuetifyEdit",
			"vuetify/views/edit.vue.tmpl",
			filename+"Edit.vue",
			data,
			false,
		)

		// edit, err := util.ExecuteTemplate("vuetify_edit.vue.tmpl", data)
		// if err == nil {
		// 	work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyEdit", Code: edit, Filename: filename + "Edit.vue"}
		// } else {
		// 	work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyEdit", Error: err, GeneratedHeaderFormat: "<-- %s -->"}
		// }
		// }(entity)
	}

	menuEntities := func() []util.Entity {
		var items []util.Entity

		for i := range entities {
			if !entities[i].Vuetify.NoGenerate {
				items = append(items, entities[i])
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
	)

	widgets := []string{
		"edit/gIcon.vue",
		"edit/gImagefield.vue",
		"edit/gMap.vue",
		"edit/gPublished.vue",
		"edit/gSelect.vue",
		"edit/gTextarea.vue",
		"edit/gTextfield.vue",
		"edit/gTime.vue",
		"edit/gToggle.vue",
		"list/gSelect.vue",
		"list/gTime.vue",
		"list/gToggle.vue",
	}

	for _, widget := range widgets {
		output.GenerateAndSave("Vuetify", "vuetify/widgets/"+widget+".tmpl", path+"/widgets/"+widget, nil, false)
	}

	// components
	output.GenerateAndSave("Vuetify", "vuetify/js/components-registration.js.tmpl", path+"/components-registration.js", struct {
		Widgets []string
	}{Widgets: widgets}, false)

	// output.GenerateAndSave("Vuetify", "vuetify/store/index.js.tmpl", path+"/store/index.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/actions.js.tmpl", path+"/store/actions.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/getters.js.tmpl", path+"/store/getters.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/mutations.js.tmpl", path+"/store/mutations.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/types.js.tmpl", path+"/store/types.js", nil, false)

	work.Waitgroup.Done()
}
