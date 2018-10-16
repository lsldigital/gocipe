package vuetify

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// Generate returns generated vuetify components
func Generate(out output.Output, r *util.Recipe) {
	if !r.Vuetify.Generate {
		// work.Waitgroup.Done()
		return
	}

	path := util.WorkingDir + "/web/" + r.Vuetify.App + "/src/gocipe"

	var forms []string
	for _, entity := range r.Entities {
		if entity.Vuetify.NoGenerate {
			continue
		}

		// go func(entity util.Entity) {
		var (
			data struct {
				Entity   util.Entity
				Entities []util.Entity
			}
		)

		data.Entity = entity
		data.Entities = r.Entities

		filename := path + "/forms/" + inflection.Plural(data.Entity.Name)

		out.GenerateAndOverwrite("VuetifyList", "vuetify/forms/list.vue.tmpl", filename+"List.vue", data)
		forms = append(forms, inflection.Plural(data.Entity.Name)+"List")

		out.GenerateAndOverwrite("VuetifyEdit", "vuetify/forms/edit.vue.tmpl", filename+"Edit.vue", data)
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

		for _, entity := range r.Entities {
			if !entity.Vuetify.NoGenerate {
				items = append(items, entity)
			}
		}

		return items
	}()

	out.GenerateAndOverwrite("Vuetify", "vuetify/js/routes.js.tmpl", path+"/routes.js", struct {
		Entities []util.Entity
	}{menuEntities})

	widgets := map[string]string{
		"EditWidgetIcon":       "edit/Icon.vue",
		"EditWidgetImagefield": "edit/Imagefield.vue",
		"EditWidgetMap":        "edit/Map.vue",
		"EditWidgetStatus":     "edit/Status.vue",
		"EditWidgetSelect":     "edit/Select.vue",
		"EditWidgetSelectRel":  "edit/SelectRel.vue",
		"EditWidgetTextarea":   "edit/Textarea.vue",
		"EditWidgetTextfield":  "edit/Textfield.vue",
		"EditWidgetTime":       "edit/Time.vue",
		"EditWidgetToggle":     "edit/Toggle.vue",
		"ListWidgetSelect":     "list/Select.vue",
		"ListWidgetTime":       "list/Time.vue",
		"ListWidgetToggle":     "list/Toggle.vue",
	}

	for _, file := range widgets {
		out.GenerateAndOverwrite("Vuetify", "vuetify/widgets/"+file+".tmpl", path+"/widgets/"+file, nil)
	}

	// components
	out.GenerateAndOverwrite("Vuetify", "vuetify/js/components-registration.js.tmpl", path+"/components-registration.js", struct {
		Widgets map[string]string
		Forms   []string
	}{Widgets: widgets, Forms: forms})

	// output.GenerateAndSave("Vuetify", "vuetify/store/index.js.tmpl", path+"/store/index.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/actions.js.tmpl", path+"/store/actions.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/getters.js.tmpl", path+"/store/getters.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/mutations.js.tmpl", path+"/store/mutations.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/types.js.tmpl", path+"/store/types.js", nil, false)
}
