package generators

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/jinzhu/inflection"
)

//StructureInfo represents a target structure in a file to be used for generation
type StructureInfo struct {
	Package   string
	Name      string
	TableName string
	Fields    []FieldInfo
}

//FieldInfo represents information about a field
//
// Example of fields (in a struct)
//
//    type Struct struct {
//        ID       *int64  `json:"id"       field.name:"id"        field.type:"serial"`
//        Authcode *string `json:"-"        field.name:"auth_code" field.type:"varchar(128)"`
//        Alias    *string `json:"alias"    field.name:"alias"     field.type:"varchar(32)"`
//        Name     *string `json:"name"     field.name:"name"      field.type:"varchar(255)"`
//        Callback *string `json:"callback" field.name:"callback"  field.type:"varchar(255)"`
//        Status   *string `json:"status"   field.name:"status"    field.type:"char(1)"`
//    }
//
type FieldInfo struct {
	Name       string // field.name
	Property   string // GO struct fields (ID, Authcode, ...)
	Type       string // GO basic value types (int64, string, ...) or custom types
	DBType     string // field.type
	Nullable   bool   // field.nullable
	Default    string // field.default
	Filterable bool   //field.filterable
	Widget     *WidgetInfo
	ManyMany   *ManyManyInfo
}

//WidgetInfo represents structure tags for widget
type WidgetInfo struct {
	Label   string
	Type    string
	Options []string
}

// ManyManyInfo represents structure for a field with a many-many relationship
type ManyManyInfo struct {
	ThisID     string
	PivotTable string
	ThatID     string
}

//NewStructureInfo process a go file to extract structure information
func NewStructureInfo(filename string, structure string) (*StructureInfo, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %s\n", err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)

	if err != nil {
		log.Fatalf("Failed to parse file: %s\n", err)
	}

	for _, d := range node.Decls {
		if decl, ok := d.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if typ, ok := spec.(*ast.TypeSpec); ok && typ.Name.Name == structure {
					return processStructure(node.Name.Name, string(src), typ)
				}
			}
		}
	}

	return nil, errors.New("Could not find structure: " + structure)
}

func processStructure(pkg string, src string, typeSpec *ast.TypeSpec) (*StructureInfo, error) {
	structInfo := new(StructureInfo)
	structInfo.Name = typeSpec.Name.Name
	structInfo.TableName = inflection.Plural(strings.ToLower(structInfo.Name))
	structInfo.Package = pkg
	structInfo.Fields = []FieldInfo{}

	if struc, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range struc.Fields.List {
			var (
				info FieldInfo
				tags reflect.StructTag
			)

			if len(field.Names) == 0 {
				continue
			}

			info.Property = field.Names[0].Name

			info.Type = strings.TrimLeft(src[field.Type.Pos()-1:field.Type.End()-1], "*")
			if field.Tag != nil && field.Tag.Value == "" {
				return nil, fmt.Errorf("structure tags not found in: %s", structInfo.Name)
			}

			tags = reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			if val, ok := tags.Lookup("field.name"); ok {
				info.Name = val
			}

			if val, ok := tags.Lookup("field.type"); ok {
				info.DBType = val
			}

			// "true" (nullable) or "false" (not null). Default: "false"
			if val, ok := tags.Lookup("field.nullable"); ok && val == "true" {
				info.Nullable = true
			}

			// "true" (filterable) or "false" (not filterable). Default: "true"
			if val, ok := tags.Lookup("field.filterable"); !ok || val != "false" {
				info.Filterable = true
			}

			if val, ok := tags.Lookup("field.default"); ok && val != "" {
				info.Default = val
			}

			if val, ok := tags.Lookup("widget"); ok && val != "" {
				w, e := ParseWidgetInfo(val)
				if e == nil {
					info.Widget = w
				} else {
					return nil, e
				}
			}

			if val, ok := tags.Lookup("field.mmany"); ok && val != "" {
				m, e := ParseManyManyInfo(val)
				if e == nil {
					info.ManyMany = m
				} else {
					return nil, e
				}
			}

			structInfo.Fields = append(structInfo.Fields, info)
		}

		return structInfo, nil
	}

	return nil, errors.New("type " + structInfo.Name + " is not a structure type.")
}

func (structInfo *StructureInfo) String() string {
	output := "\n"
	output += "Structure Name: " + structInfo.Name + "\n\n"
	output += fmt.Sprintf("\t%10s\t%10s\n", "Name:", "Type:")
	output += fmt.Sprintf("\t%10s\t%10s\n", "-----", "-----")
	for _, fieldInfo := range structInfo.Fields {
		output += fmt.Sprintf("\t%10s\t%10s\n", fieldInfo.Name, fieldInfo.Type)
	}

	return output
}

//ParseWidgetInfo returns WidgetInfo by parsing a string; format being Label#Type#Options
func ParseWidgetInfo(value string) (*WidgetInfo, error) {
	var (
		widgetInfo WidgetInfo
		fields     []string
		lenf       int
	)

	fields = strings.SplitN(value, "#", 3)
	lenf = len(fields)

	if lenf != 2 && lenf != 3 {
		return nil, errors.New("invalid format for widget tag: " + value)
	}

	widgetInfo.Label = fields[0]
	widgetInfo.Type = fields[1]

	if lenf == 3 {
		widgetInfo.Options = strings.Split(fields[2], ";")
	}

	return &widgetInfo, nil
}

// ParseManyManyInfo returns ManyManyInfo by parsing a string; format being ThisID#PivotTable#ThatID
func ParseManyManyInfo(value string) (*ManyManyInfo, error) {
	var (
		manyMany ManyManyInfo
		fields   []string
	)

	fields = strings.Split(value, "#")
	if len(fields) != 3 {
		return nil, errors.New("invalid format for many many tag: " + value)
	}

	manyMany.ThisID = fields[0]
	manyMany.PivotTable = fields[1]
	manyMany.ThatID = fields[2]

	return &manyMany, nil
}
