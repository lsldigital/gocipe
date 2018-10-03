package util

import "errors"

var (
	//ErrorRecipeEntityNotFound indicates that an entity was not found during lookup request
	ErrorRecipeEntityNotFound = errors.New("recipe does not contain requested entity")
)

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

		for j := range e.References {
			c := &e.References[j]
			c.init()
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
