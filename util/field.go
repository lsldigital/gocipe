package util

import (
	"encoding/json"
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

	//ErrorFieldTypeInvalid indicates field type is not supported
	ErrorFieldTypeInvalid = errors.New("field type is not valid")
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
	switch f.Type { //TODO SQL-dialect sensitive
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

	switch f.Type {
	default:
		return ErrorFieldTypeInvalid
	case fieldTypeBool, fieldTypeInt, fieldTypeStr, fieldTypeTime:
		//all good
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

//ProtoDefinition returns proto definition for this field
func (f *Field) ProtoDefinition(index *int) string {
	var protoType string

	switch f.Type {
	case fieldTypeBool:
		protoType = "bool"
	case fieldTypeInt:
		protoType = "int64"
	case fieldTypeStr:
		protoType = "string"
	case fieldTypeTime:
		protoType = "google.protobuf.Timestamp"
	}

	definition := fmt.Sprintf(`%s %s = %d;`, protoType, f.Name, *index)
	*index++

	return definition
}

//GetBefore returns code executed before a statement is executed (used in crud)
func (f *Field) GetBefore(op string) []string {
	var before []string

	if f.Name == "CreatedAt" {
		if op == "insert" {
			before = append(before, fmt.Sprintf(`entity.%s = ptypes.TimestampNow()`, f.Name))
		}
	} else if f.Name == "UpdatedAt" {
		before = append(before, fmt.Sprintf(`entity.%s = ptypes.TimestampNow()`, f.Name))
	}

	//some fields require preprocessing
	//they will be assigned to a variable, use that instead of the property name
	if f.Type == "time" {
		switch op {
		case "get", "list":
			before = append(before, fmt.Sprintf(`var %s time.Time`, strings.ToLower(f.Name)))
		case "insert", "merge", "update":
			before = append(before, fmt.Sprintf(`%s, _ := ptypes.Timestamp(entity.%s)`, strings.ToLower(f.Name), f.Name))
		}
	}

	return before
}

//GetAfter returns code executed after a statement is executed (used in crud)
func (f *Field) GetAfter(op string) []string {
	var after []string

	//some fields require preprocessing
	//they will be assigned to a variable, use that instead of the property name
	if f.Type == "time" {
		switch op {
		case "get", "list":
			after = append(after, fmt.Sprintf(`entity.%s, _ = ptypes.TimestampProto(%s)`, f.Name, strings.ToLower(f.Name)))
		}
	}

	return after
}

//GetEditWidgetOptionsJSON returns code executed after a statement is executed (used in crud)
func (f *Field) GetEditWidgetOptionsJSON() string {
	jsob, err := json.Marshal(f.EditWidget.Options)
	if err != nil {
		return ""
	}

	return string(jsob)
}
