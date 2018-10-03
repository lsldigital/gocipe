package util

import (
	"errors"
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

// Field describes a field contained in an entity
type Field struct {
	// Name is the name of the property
	Name string `json:"name"`

	// Label is the label for the field
	Label string `json:"label"`

	// Type is the data type of the property
	Type string `json:"type"`

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
