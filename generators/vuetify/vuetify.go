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

	path := util.WorkingDir + "/web/" + recipe.Vuetify.App + "/src/"

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

		filename := path + "gocipe/forms/" + inflection.Plural(data.Entity.Name)

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
		path+"gocipe/routes.js",
		struct {
			Entities []util.Entity
		}{menuEntities},
		false,
	)

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
		output.GenerateAndSave("Vuetify", "vuetify/widgets/"+file+".tmpl", path+"gocipe/widgets/"+file, nil, false)
	}
	// components
	output.GenerateAndSave("Vuetify", "vuetify/js/components-registration.js.tmpl", path+"gocipe/components-registration.js", struct {
		Widgets map[string]string
		Forms   []string
	}{Widgets: widgets, Forms: forms}, false)

	staticfiles := map[string]string{
		"shared-ui/AppFooter.vue":         "shared-ui/AppFooter.vue",
		"shared-ui/AppNavigation.vue":     "shared-ui/AppNavigation.vue",
		"shared-ui/AppToolbar.vue":        "shared-ui/AppToolbar.vue",
		"shared-ui/NotFound.vue":          "shared-ui/NotFound.vue",
		"shared-ui/PageHome.vue":          "shared-ui/PageHome.vue",
		"shared-ui/Authenticated.vue":     "shared-ui/Authenticated.vue",
		"shared-ui/Login.vue":             "shared-ui/Login.vue",
		"store/modules/auth/index.js":     "store/modules/auth/index.js",
		"store/modules/auth/getters.js":   "store/modules/auth/getters.js",
		"store/modules/auth/actions.js":   "store/modules/auth/actions.js",
		"store/modules/auth/mutations.js": "store/modules/auth/mutations.js",
		"store/modules/auth/types.js":     "store/modules/auth/types.js",
	}

	for src, target := range staticfiles {
		output.GenerateAndSave("Vuetify", "vuetify/"+src+".tmpl", path+"gocipe/"+target, nil, false)
	}

	output.GenerateAndSave("Vuetify", "vuetify/js/router.js.tmpl", path+"router.js", nil, false)
	output.GenerateAndSave("Vuetify", "vuetify/js/store.js.tmpl", path+"store.js", nil, false)
	output.GenerateAndSave("Vuetify", "vuetify/shared-ui/App.vue.tmpl", path+"App.vue", nil, false)

	// output.GenerateAndSave("Vuetify", "vuetify/store/index.js.tmpl", path+"gocipe/store/index.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/actions.js.tmpl", path+"gocipe/store/actions.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/getters.js.tmpl", path+"gocipe/store/getters.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/mutations.js.tmpl", path+"gocipe/store/mutations.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/types.js.tmpl", path+"gocipe/store/types.js", nil, false)

	work.Waitgroup.Done()
}
