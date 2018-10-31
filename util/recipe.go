package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
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

	// Crud describes options for Crud generation
	Crud CrudOpts `json:"crud"`

	// Admin describes options for Admin service generation
	Admin AdminOpts `json:"admin"`

	// Vuetify describes options for Vuetify generation
	Vuetify VuetifyOpts `json:"vuetify"`

	// Decks describes options for Decks generation
	Decks DecksOpts `json:"decks"`

	// Entities lists entities to be generated
	Entities []Entity `json:"entities"`

	//entities is a map for random access of entities contained in entity
	entities map[string]*Entity
}

// LoadRecipe loads gocipe config file and returns it as Recipe
func LoadRecipe() (*Recipe, error) {
	var recipe Recipe

	recipePath, err := GetAbsPath("gocipe.json")
	if err != nil {
		return nil, err
	}

	if !FileExists(recipePath) {
		return nil, fmt.Errorf("file not found: %s", recipePath)
	}

	recipeContent, err := ioutil.ReadFile(recipePath)

	if err != nil {
		return nil, err
	}
	//ToDo
	// output.Log("%x", sha256.Sum256([]byte(recipeContent)))

	// output.Inject(
	// 	recipePath,
	// )

	err = json.Unmarshal(recipeContent, &recipe)
	if err != nil {
		return nil, fmt.Errorf("recipe decoding failed: %s", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	WorkingDir = wd
	AppImportPath = strings.TrimPrefix(wd, os.Getenv("GOPATH")+"/src/")
	AppName = path.Base(AppImportPath)
	if AppName == "." {
		AppName = "app"
	}
	os.Getenv("GOPATH")

	recipe.init()

	return &recipe, nil
}

func (r *Recipe) init() {
	r.entities = make(map[string]*Entity)

	// Add default entities
	if r.Decks.Generate {
		r.Entities = append(r.Entities, card, cardSchedule)
	}

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
	var (
		permissions []Permission
		index       int
	)

	ops := []string{"Create", "Edit", "View", "List", "Delete", "Lookup"}
	utf8list := genUTF8List(len(r.Entities) * len(ops) * 3)
	for _, e := range r.Entities {
		if !e.Admin.Generate {
			continue
		}
		for _, o := range ops {
			permissions = append(permissions, Permission{
				Name: "Perm" + o + e.Name,
				Code: utf8list[index],
			})
			index++
			if o != "Create" {
				permissions = append(permissions, Permission{
					Name: "Perm" + o + e.Name + "Any",
					Code: utf8list[index],
				})
				index++
			}
		}
	}

	return permissions
}

//HasFileFields returns true if at least an entity has file fields
func (r *Recipe) HasFileFields() bool {
	for _, e := range r.Entities {
		if e.HasFileFields() {
			return true
		}
	}
	return false
}

//HasContentFileUpload returns true if at least an entity has file upload option for content builder
func (r *Recipe) HasContentFileUpload() bool {
	for _, e := range r.Entities {
		if e.ContentBuilder.Generate && e.ContentBuilder.AllowUpload {
			return true
		}
	}
	return false
}

//HasAuth returns true if at least an entity requires Auth
func (r *Recipe) HasAuth() bool {
	for _, e := range r.Entities {
		if e.Admin.Auth.Generate {
			return true
		}
	}
	return false
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
