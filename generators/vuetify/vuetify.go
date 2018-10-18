package vuetify

import (
	"path"
	"path/filepath"

	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// Generate returns generated vuetify components
func Generate(out *output.Output, r *util.Recipe) {
	if !r.Vuetify.Generate {
		// work.Waitgroup.Done()
		return
	}

	var (
		forms        []string
		menuEntities []util.Entity
	)
	dstPath := path.Join(util.WorkingDir, "/web/", r.Vuetify.App, "/src/gocipe")

	for _, entity := range r.Entities {
		if entity.Vuetify.NoGenerate {
			continue
		}

		if !entity.Vuetify.NoGenerate {
			menuEntities = append(menuEntities, entity)
		}

		// go func(entity util.Entity) {
		data := struct {
			Entity util.Entity
		}{entity}

		filePath := path.Join(dstPath, "/forms/")
		fileName := inflection.Plural(entity.Name)

		out.GenerateAndOverwrite("Vuetify List", "vuetify/forms/list.vue.tmpl", filepath.Join(filePath, fileName+"List.vue"), data)
		forms = append(forms, inflection.Plural(data.Entity.Name)+"List")

		out.GenerateAndOverwrite("Vuetify Edit", "vuetify/forms/edit.vue.tmpl", filepath.Join(filePath, fileName+"Edit.vue"), data)
		forms = append(forms, inflection.Plural(data.Entity.Name)+"Edit")

		// edit, err := util.ExecuteTemplate("vuetify_edit.vue.tmpl", data)
		// if err == nil {
		// 	work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyEdit", Code: edit, filePath: filePath + "Edit.vue"}
		// } else {
		// 	work.Done <- util.GeneratedCode{Generator: "GenerateVuetifyEdit", Error: err, GeneratedHeaderFormat: "<-- %s -->"}
		// }
		// }(entity)
	}

	out.GenerateAndOverwrite("Vuetify", "vuetify/js/routes.js.tmpl", filepath.Join(dstPath, "/routes.js"), struct {
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
		out.GenerateAndOverwrite("Vuetify Widgets", filepath.Join("vuetify/widgets/", file+".tmpl"), filepath.Join(dstPath, "/widgets/", file), nil)
	}

	// components
	out.GenerateAndOverwrite("Vuetify Registration", "vuetify/js/components-registration.js.tmpl", filepath.Join(dstPath, "/components-registration.js"), struct {
		Widgets map[string]string
		Forms   []string
	}{Widgets: widgets, Forms: forms})

	// output.GenerateAndSave("Vuetify", "vuetify/store/index.js.tmpl", path+"/store/index.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/actions.js.tmpl", path+"/store/actions.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/getters.js.tmpl", path+"/store/getters.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/mutations.js.tmpl", path+"/store/mutations.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/types.js.tmpl", path+"/store/types.js", nil, false)
}
