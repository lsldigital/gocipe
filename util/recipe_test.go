package util

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"flag"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const testDataDir = "./test-fixtures/"

var update = flag.Bool("update", false, "update golden files")

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
			args:    args{recipeFile: "gocipe_malformed"},
			wantErr: true,
		},
		{
			name:    "Duplicate entries recipe file",
			args:    args{recipeFile: "gocipe_duplicate"},
			wantErr: false,
		},
		{
			name:    "Sample recipe file",
			args:    args{recipeFile: "gocipe"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := LoadRecipe(filepath.Join(testDataDir, tt.args.recipeFile+".json"))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRecipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && (actual == nil) {
				return
			}

			golden := filepath.Join(testDataDir, tt.args.recipeFile+".golden")
			if *update {
				testWriteGob(t, golden, actual)
			}

			expected := new(Recipe)
			testReadGob(t, golden, expected)
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("LoadRecipe() actual = %+v, \n\n\n want %+v", actual, expected)
				return
			}
		})
	}
}

func testRecipeToByte(t *testing.T, recipe *Recipe) []byte {
	t.Helper()

	buf := &bytes.Buffer{}
	log.Println(recipe)
	err := binary.Write(buf, binary.LittleEndian, *recipe)
	if err != nil {
		t.Fatalf("Converting recipe to bytes err: %s", err)
	}
	return buf.Bytes()
}

func testWriteGob(t *testing.T, filePath string, object interface{}) {
	t.Helper()

	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()

	if err != nil {
		t.Fatalf("Writing gob err: %s", err)
	}
}

func testReadGob(t *testing.T, filePath string, object interface{}) {
	t.Helper()

	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()

	if err != nil {
		t.Fatalf("Reading gob err: %s", err)
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
