package util

import (
	"reflect"
	"testing"
)

func TestLoadRecipe(t *testing.T) {
	type args struct {
		recipeFile string
	}
	tests := []struct {
		name    string
		args    args
		want    *Recipe
		wantErr bool
	}{
		{
			name:    "No recipe file",
			args:    args{recipeFile: ""},
			wantErr: true,
		},
		{
			name:    "Malformed recipe file",
			args:    args{recipeFile: "gocipe_test_malformed.json"},
			wantErr: true,
		},
		{
			name:    "Duplicate recipe entries file",
			args:    args{recipeFile: "gocipe_test_duplicate.json"},
			wantErr: false,
		},
		{
			name:    "Sample recipe file",
			args:    args{recipeFile: "gocipe_test.json"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadRecipe(tt.args.recipeFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRecipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("LoadRecipe() = %v, want recipe object", got)
			}
		})
	}
}

func TestRecipe_init(t *testing.T) {
	type fields struct {
		ImportPath string
		Bootstrap  BootstrapOpts
		Crud       CrudOpts
		Admin      AdminOpts
		Vuetify    VuetifyOpts
		Decks      DecksOpts
		Entities   []Entity
		entities   map[string]*Entity
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recipe{
				ImportPath: tt.fields.ImportPath,
				Bootstrap:  tt.fields.Bootstrap,
				Crud:       tt.fields.Crud,
				Admin:      tt.fields.Admin,
				Vuetify:    tt.fields.Vuetify,
				Decks:      tt.fields.Decks,
				Entities:   tt.fields.Entities,
				entities:   tt.fields.entities,
			}
			r.init()
		})
	}
}

func TestRecipe_GetEntity(t *testing.T) {
	type fields struct {
		ImportPath string
		Bootstrap  BootstrapOpts
		Crud       CrudOpts
		Admin      AdminOpts
		Vuetify    VuetifyOpts
		Decks      DecksOpts
		Entities   []Entity
		entities   map[string]*Entity
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Entity
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recipe{
				ImportPath: tt.fields.ImportPath,
				Bootstrap:  tt.fields.Bootstrap,
				Crud:       tt.fields.Crud,
				Admin:      tt.fields.Admin,
				Vuetify:    tt.fields.Vuetify,
				Decks:      tt.fields.Decks,
				Entities:   tt.fields.Entities,
				entities:   tt.fields.entities,
			}
			got, err := r.GetEntity(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Recipe.GetEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Recipe.GetEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecipe_Validate(t *testing.T) {
	type fields struct {
		ImportPath string
		Bootstrap  BootstrapOpts
		Crud       CrudOpts
		Admin      AdminOpts
		Vuetify    VuetifyOpts
		Decks      DecksOpts
		Entities   []Entity
		entities   map[string]*Entity
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recipe{
				ImportPath: tt.fields.ImportPath,
				Bootstrap:  tt.fields.Bootstrap,
				Crud:       tt.fields.Crud,
				Admin:      tt.fields.Admin,
				Vuetify:    tt.fields.Vuetify,
				Decks:      tt.fields.Decks,
				Entities:   tt.fields.Entities,
				entities:   tt.fields.entities,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Recipe.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRecipe_GetPermissions(t *testing.T) {
	type fields struct {
		ImportPath string
		Bootstrap  BootstrapOpts
		Crud       CrudOpts
		Admin      AdminOpts
		Vuetify    VuetifyOpts
		Decks      DecksOpts
		Entities   []Entity
		entities   map[string]*Entity
	}
	tests := []struct {
		name   string
		fields fields
		want   []Permission
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recipe{
				ImportPath: tt.fields.ImportPath,
				Bootstrap:  tt.fields.Bootstrap,
				Crud:       tt.fields.Crud,
				Admin:      tt.fields.Admin,
				Vuetify:    tt.fields.Vuetify,
				Decks:      tt.fields.Decks,
				Entities:   tt.fields.Entities,
				entities:   tt.fields.entities,
			}
			if got := r.GetPermissions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Recipe.GetPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecipe_HasFileFields(t *testing.T) {
	type fields struct {
		ImportPath string
		Bootstrap  BootstrapOpts
		Crud       CrudOpts
		Admin      AdminOpts
		Vuetify    VuetifyOpts
		Decks      DecksOpts
		Entities   []Entity
		entities   map[string]*Entity
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recipe{
				ImportPath: tt.fields.ImportPath,
				Bootstrap:  tt.fields.Bootstrap,
				Crud:       tt.fields.Crud,
				Admin:      tt.fields.Admin,
				Vuetify:    tt.fields.Vuetify,
				Decks:      tt.fields.Decks,
				Entities:   tt.fields.Entities,
				entities:   tt.fields.entities,
			}
			if got := r.HasFileFields(); got != tt.want {
				t.Errorf("Recipe.HasFileFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecipe_HasContentFileUpload(t *testing.T) {
	type fields struct {
		ImportPath string
		Bootstrap  BootstrapOpts
		Crud       CrudOpts
		Admin      AdminOpts
		Vuetify    VuetifyOpts
		Decks      DecksOpts
		Entities   []Entity
		entities   map[string]*Entity
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recipe{
				ImportPath: tt.fields.ImportPath,
				Bootstrap:  tt.fields.Bootstrap,
				Crud:       tt.fields.Crud,
				Admin:      tt.fields.Admin,
				Vuetify:    tt.fields.Vuetify,
				Decks:      tt.fields.Decks,
				Entities:   tt.fields.Entities,
				entities:   tt.fields.entities,
			}
			if got := r.HasContentFileUpload(); got != tt.want {
				t.Errorf("Recipe.HasContentFileUpload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecipe_HasAuth(t *testing.T) {
	type fields struct {
		ImportPath string
		Bootstrap  BootstrapOpts
		Crud       CrudOpts
		Admin      AdminOpts
		Vuetify    VuetifyOpts
		Decks      DecksOpts
		Entities   []Entity
		entities   map[string]*Entity
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recipe{
				ImportPath: tt.fields.ImportPath,
				Bootstrap:  tt.fields.Bootstrap,
				Crud:       tt.fields.Crud,
				Admin:      tt.fields.Admin,
				Vuetify:    tt.fields.Vuetify,
				Decks:      tt.fields.Decks,
				Entities:   tt.fields.Entities,
				entities:   tt.fields.entities,
			}
			if got := r.HasAuth(); got != tt.want {
				t.Errorf("Recipe.HasAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genUTF8List(t *testing.T) {
	type args struct {
		limit int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genUTF8List(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genUTF8List() = %v, want %v", got, tt.want)
			}
		})
	}
}
