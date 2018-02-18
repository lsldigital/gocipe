package generators

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

//StructureInfo represents a target structure in a file to be used for generation
type StructureInfo struct {
	Package string
	Name    string
	Fields  []FieldInfo
}

//FieldInfo represents information about a field
type FieldInfo struct {
	Name     string
	Type     string
	Comments string
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
					return ProcessStructure(node.Name.Name, string(src), typ)
				}
			}
		}
	}

	return nil, errors.New("Could not find structure: " + structure)
}

//ProcessStructure creates an returns a structure given a TypeSpec
func ProcessStructure(pkg string, src string, typeSpec *ast.TypeSpec) (*StructureInfo, error) {
	structInfo := new(StructureInfo)
	structInfo.Name = typeSpec.Name.Name
	structInfo.Package = pkg
	structInfo.Fields = make([]FieldInfo, 1)

	if struc, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range struc.Fields.List {
			if len(field.Names) == 0 {
				continue
			}
			fieldtype := src[field.Type.Pos()-1 : field.Type.End()-1]
			structInfo.Fields = append(structInfo.Fields, FieldInfo{Name: field.Names[0].Name, Type: fieldtype, Comments: field.Comment.Text()})
		}

		return structInfo, nil
	}

	return nil, errors.New("Type " + structInfo.Name + " is not a structure type.")
}

func (structInfo *StructureInfo) String() string {
	output := "\n"
	output += "Structure Name: " + structInfo.Name + "\n\n"
	output += fmt.Sprintf("\t%10s\t%10s\t%s\n", "Name:", "Type:", "Comments:")
	output += fmt.Sprintf("\t%10s\t%10s\t%s\n", "-----", "-----", "---------")
	for _, fieldInfo := range structInfo.Fields {
		output += fmt.Sprintf("\t%10s\t%10s\t%s\n", fieldInfo.Name, fieldInfo.Type, fieldInfo.Comments)
	}

	return output
}
