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

	var forms []string
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

		filename := path + "/forms/" + inflection.Plural(data.Entity.Name)

		output.GenerateAndSave(
			"VuetifyList",
			"vuetify/forms/list.vue.tmpl",
			filename+"List.vue",
			data,
			false,
		)
		forms = append(forms, inflection.Plural(data.Entity.Name)+"List")

		output.GenerateAndSave(
			"VuetifyEdit",
			"vuetify/forms/edit.vue.tmpl",
			filename+"Edit.vue",
			data,
			false,
		)
		forms = append(forms, inflection.Plural(data.Entity.Name)+"Edit")

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

	widgets := map[string]string{
		"edit/gIcon.vue":       "gIcon",
		"edit/gImagefield.vue": "gImagefield",
		"edit/gMap.vue":        "gMap",
		"edit/gPublished.vue":  "gPublished",
		"edit/gSelect.vue":     "gSelect",
		"edit/gTextarea.vue":   "gTextarea",
		"edit/gTextfield.vue":  "gTextfield",
		"edit/gTime.vue":       "gTime",
		"edit/gToggle.vue":     "gToggle",
		"list/gSelect.vue":     "gSelect",
		"list/gTime.vue":       "gTime",
		"list/gToggle.vue":     "gToggle",
	}

	for file := range widgets {
		output.GenerateAndSave("Vuetify", "vuetify/widgets/"+file+".tmpl", path+"/widgets/"+file, nil, false)
	}

	// components
	output.GenerateAndSave("Vuetify", "vuetify/js/components-registration.js.tmpl", path+"/components-registration.js", struct {
		Widgets map[string]string
		Forms   []string
	}{Widgets: widgets, Forms: forms}, false)

	// output.GenerateAndSave("Vuetify", "vuetify/store/index.js.tmpl", path+"/store/index.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/actions.js.tmpl", path+"/store/actions.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/getters.js.tmpl", path+"/store/getters.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/mutations.js.tmpl", path+"/store/mutations.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/types.js.tmpl", path+"/store/types.js", nil, false)

	work.Waitgroup.Done()
}
