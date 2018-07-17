package bread

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// GenerateBREAD returns generated code for a BREAD service - Browse, Read, Edit, Add & Delete
func GenerateBREAD(work util.GenerationWork, entities []util.Entity) error {
	var ents []util.Entity

	//3 jobs to be waited upon for: service_bread.proto, service_bread.go and service_bread_hooks.go
	work.Waitgroup.Add(3)

	for i := range entities {
		if entities[i].Bread.Generate {
			ents = append(ents, entities[i])
		}
	}

	code, err := util.ExecuteTemplate("service_bread.go.tmpl", struct {
		Entities []util.Entity
	}{ents})

	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBread", Code: code, Filename: "services/bread/service_bread.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateBread", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}

	for _, entity := range entities {
		switch true {
		case
			entity.Bread.Hooks.PreCreate,
			entity.Bread.Hooks.PostCreate,
			entity.Bread.Hooks.PreRead,
			entity.Bread.Hooks.PostRead,
			entity.Bread.Hooks.PreList,
			entity.Bread.Hooks.PostList,
			entity.Bread.Hooks.PreUpdate,
			entity.Bread.Hooks.PostUpdate,
			entity.Bread.Hooks.PreDelete,
			entity.Bread.Hooks.PostDelete:
			hooks, err := util.ExecuteTemplate("service_bread_hooks.go.tmpl", struct {
				Entity util.Entity
			}{entity})

			if err == nil {
				work.Done <- util.GeneratedCode{
					Generator: "GenerateBreadHooks:" + entity.Name,
					Code:      hooks,
					Filename: fmt.Sprintf(
						"services/bread/service_bread_%s_hooks.gocipe.go",
						strings.ToLower(entity.Name),
					),
					NoOverwrite: true,
				}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateBreadHooks", Error: err}
			}
		}
	}

	proto, err := util.ExecuteTemplate("service_bread.proto.tmpl", struct {
		Entities []util.Entity
	}{ents})

	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBreadProto", Code: proto, Filename: "proto/service_bread.proto"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateBreadProto", Error: err}
	}

	work.Waitgroup.Done()
	return nil
}
