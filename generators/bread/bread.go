package bread

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// Generate returns generated code for a BREAD service - Browse, Read, Edit, Add & Delete
func Generate(work util.GenerationWork, entities map[string]util.Entity) error {
	var ents []util.Entity

	//2 jobs to be waited upon for: service_bread.proto and service_bread.go
	work.Waitgroup.Add(2)

	for i := range entities {
		if entities[i].Bread.Generate {
			ents = append(ents, entities[i])
		}
	}

	code, err := util.ExecuteTemplate("bread/service_bread.go.tmpl", struct {
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
			hooks, err := util.ExecuteTemplate("bread/service_bread_hooks.go.tmpl", struct {
				Entity     util.Entity
				PreRead    bool
				PostRead   bool
				PreList    bool
				PostList   bool
				PreCreate  bool
				PostCreate bool
				PreUpdate  bool
				PostUpdate bool
				PreDelete  bool
				PostDelete bool
			}{
				Entity:     entity,
				PreRead:    entity.Bread.Hooks.PreRead,
				PostRead:   entity.Bread.Hooks.PostRead,
				PreList:    entity.Bread.Hooks.PreList,
				PostList:   entity.Bread.Hooks.PostList,
				PreCreate:  entity.Bread.Hooks.PreCreate,
				PostCreate: entity.Bread.Hooks.PostCreate,
				PreUpdate:  entity.Bread.Hooks.PreUpdate,
				PostUpdate: entity.Bread.Hooks.PostUpdate,
				PreDelete:  entity.Bread.Hooks.PreDelete,
				PostDelete: entity.Bread.Hooks.PostDelete,
			})

			work.Waitgroup.Add(1)

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

	proto, err := util.ExecuteTemplate("bread/service_bread.proto.tmpl", struct {
		Entities          []util.Entity
		ProjectImportPath string
	}{ents, util.ProjectImportPath})

	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBreadProto", Code: proto, Filename: "proto/service_bread.proto"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateBreadProto", Error: err}
	}

	work.Waitgroup.Done()
	return nil
}
