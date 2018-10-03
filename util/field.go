package util

import (
	"errors"
	"fmt"
	"strings"
)

const (
	fieldTypeBool = "bool"
	fieldTypeInt  = "int64"
	fieldTypeStr  = "string"
	fieldTypeTime = "time"
)

var (
	//ErrorFieldNameReserved indicates a fieldname is reserved
	ErrorFieldNameReserved = errors.New("field name is reserved")

	//ErrorFieldNameEmpty indicates a fieldname being empty
	ErrorFieldNameEmpty = errors.New("field name is empty")
)

// FieldSchema represents schema generation information for the field
type FieldSchema struct {
	// Field is the name of the field in database
	Field string `json:"field"`

	// Type is the data type for the field in database
	Type string `json:"type"`
}

// Field describes a field contained in an entity
type Field struct {
	// Name is the name of the property
	Name string `json:"name"`

	// Label is the label for the field
	Label string `json:"label"`

	// Type is the data type of the property
	Type string `json:"type"`

	// Default provides the default value for this field in database
	Default string `json:"default"`

	// EditWidget represents widget information for the field
	EditWidget EditWidgetOpts `json:"edit_widget"`

	// ListWidget represents widget information for the field
	ListWidget ListWidgetOpts `json:"list_widget"`

	// schema represents schema information for the field
	schema FieldSchema
}

func (f *Field) init() {
	switch f.Type {
	case fieldTypeBool:
		f.schema.Type = "BOOL"
	case fieldTypeInt:
		f.schema.Type = "INTEGER"
	case fieldTypeStr:
		f.schema.Type = "TEXT"
	case fieldTypeTime:
		f.schema.Type = "TIMESTAMP"
	}

	f.schema.Field = ToSnakeCase(f.Name)
}

//Validate checks for errors in the field
func (f *Field) Validate() error {
	if f.Name == "" {
		return ErrorFieldNameEmpty
	}

	switch f.Name {
	case "Status", "ID", "UserID":
		return ErrorFieldNameReserved
	}

	return nil
}

//SchemaDefinition returns schema definition for this field as part of a table definition
func (f *Field) SchemaDefinition() string {
	var def string
	if f.Default != "" {
		def = fmt.Sprintf(` DEFAULT "%s"`, strings.Replace(f.Default, `"`, `""`, -1)) //TODO SQL-dialect sensitive
	}
	return fmt.Sprintf(`"%s" %s NOT NULL%s`, f.schema.Field, f.schema.Type, def)
}
