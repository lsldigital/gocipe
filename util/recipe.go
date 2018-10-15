package util

import (
	"errors"
	"unicode"
)

var (
	//ErrorRecipeEntityNotFound indicates that an entity was not found during lookup request
	ErrorRecipeEntityNotFound = errors.New("recipe does not contain requested entity")
)

// Permission used in GetPermissions
type Permission struct {
	Name, Code string
}

// Recipe represents a recipe to generate a project
type Recipe struct {
	// Container indicates whether or not container should be generated
	Bootstrap BootstrapOpts `json:"bootstrap"`

	// HTTP indicates whether http server code should be generated
	HTTP HTTPOpts `json:"http"`

	// Crud describes options for Crud generation
	Crud CrudOpts `json:"crud"`

	// Admin describes options for Admin service generation
	Admin AdminOpts `json:"admin"`

	// Vuetify describes options for Vuetify generation
	Vuetify VuetifyOpts `json:"vuetify"`

	// Entities lists entities to be generated
	Entities []Entity `json:"entities"`

	//entities is a map for random access of entities contained in entity
	entities map[string]*Entity
}

func (r *Recipe) init() {
	r.entities = make(map[string]*Entity)

	for i := range r.Entities {
		e := &r.Entities[i]
		e.init(r)
		r.entities[e.Name] = e
	}

	for i := range r.Entities {
		e := &r.Entities[i]
		for j := range e.Relationships { //relationships must be init after all entities are init
			p := &e.Relationships[j]
			p.init(r, e)
		}
	}
}

//GetEntity returns an entity in the recipe by name
func (r *Recipe) GetEntity(name string) (*Entity, error) {
	if e, ok := r.entities[name]; ok {
		return e, nil
	}

	return nil, ErrorRecipeEntityNotFound
}

//Validate checks the recipe for errors
func (r *Recipe) Validate() error {
	return nil
}

//GetPermissions returns a map of permissions
func (r *Recipe) GetPermissions() []Permission {
	var permissions []Permission

	ops := []string{"Create", "Edit", "View", "List", "Delete", "Lookup"}
	utf8list := genUTF8List(len(r.Entities) * len(ops) * 2)
	for _, e := range r.Entities {
		for _, o := range ops {
			for y := 0; y < len(utf8list); y++ {
				permissions = append(permissions, Permission{
					Name: "Perm" + o + e.Name,
					Code: utf8list[y],
				})
				if o != "Create" {
					y++
					permissions = append(permissions, Permission{
						Name: "Perm" + o + e.Name + "Any",
						Code: utf8list[y],
					})
				}
			}
		}
	}

	return permissions
}

func genUTF8List(limit int) []string {
	var (
		strList []string
		r       rune
	)
	for i := 0; i < limit; i++ {
		r = 'A' + rune(i)
		if unicode.IsPrint(r) &&
			!unicode.IsSymbol(r) &&
			!unicode.IsSpace(r) &&
			!unicode.IsControl(r) &&
			r != '\\' {
			strList = append(strList, string(r))
		}
	}
	return strList
}
