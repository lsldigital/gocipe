package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/generators/crud"
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// // RepoCodes that must be analysed for interface generation
// var (
// 	RepoCodes []RepoCode
// 	repoChan  = make(chan RepoCode)
// )

// // RepoCode represents generated code for a crud repo, which will be used to extract an interface
// type RepoCode struct {
// 	SourceFile string
// 	TargetFile string
// }

type entityCrud struct {
	Imports      []string
	Structure    string
	Get          string
	List         string
	DeleteSingle string
	DeleteMany   string
	Save         string
	Insert       string
	Update       string
	Merge        string
	SaveRelated  []string
	LoadRelated  []string
}

type relationship struct {
	Table        string
	ThisID       string
	ThatID       string
	PropertyName string
}

// Generate returns generated code to run an http server
func Generate(out *output.Output, r *util.Recipe) {
	var generateAny bool
	// work.Waitgroup.Add(len(entities) * 2) //2 threads per entities. for models and models_hooks

	for _, e := range r.Entities {

		if e.HasCrudHooks() {
			output.GenerateAndSave("crud:hooks", "crud/hooks.go.tmpl", fmt.Sprintf("models/%s_crud_hooks.gocipe.go", strings.ToLower(e.Name)), e, true)
		}
	}

	out.GenerateAndOverwrite("Crud Proto", "crud/models.proto.tmpl", "proto/models.proto", output.WithHeader, struct {
		Entities      []util.Entity
		AppImportPath string
	}{Entities: r.Entities, AppImportPath: util.AppImportPath})

	out.GenerateAndOverwrite("Crud Moderrors", "crud/moderrors.go.tmpl", "models/moderrors/errors.gocipe.go", output.WithHeader, struct{}{})

	//old:

	for _, entity := range r.Entities {
		generateAny = generateAny || r.Crud.Generate

		if !r.Crud.Generate {
			util.DeleteIfExists(fmt.Sprintf("models/%s_repo.gocipe.go", strings.ToLower(entity.Name)))
			util.DeleteIfExists(fmt.Sprintf("models/%s_crud_hooks.gocipe.go", strings.ToLower(entity.Name)))
			// work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUD[%s]", entity.Name), Error: util.ErrorSkip}
			// work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUDHooks[%s]", entity.Name), Error: util.ErrorSkip}
			//TODO cater for skip
			continue
		}

		if code, err := generateCrud2(r); err == nil {
			name := strings.ToLower(entity.Name)
			out.GenerateAndOverwrite("Crud Entity "+name, "crud/crud.go.tmpl", fmt.Sprintf("models/%s_repo.gocipe.go", name), output.WithHeader, code)
		} //TODO cater for error
	}

	if generateAny {
		out.GenerateAndOverwrite("crud/models.go.tmpl", "models/models.gocipe.go", output.WithHeader, struct {
			Crud     bool
			Entities map[string]util.Entity
		}{
			Crud:     crud.Generate,
			Entities: entities,
		})

	}

}

func generateCrud2(r *util.Recipe) (string, error) {
	var (
		entities []util.Postgres
		imports  []string
	)

	for _, e := range r.Entities {
		entities = append(entities, util.Postgres{Entity: e})
	}

	imports = append(imports, `uuid "github.com/satori/go.uuid"`)

	return util.ExecuteTemplate("crud/crud.go.tmpl", struct {
		Entities []util.Postgres
		Imports  []string
	}{Entities: entities, Imports: imports})
}
