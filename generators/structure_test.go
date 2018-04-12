package generators

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	structInfo := StructureInfo{

		Name:      "Persons",
		TableName: "Persons",
		Fields: []FieldInfo{
			{Name: "id", Type: "int64", Tags: `json:"id" field.name:"id" field.type:"serial"`},
			{Name: "name", Type: "string", Tags: `json:"name" field.name:"name" field.type:"varchar(255)"`},
			{Name: "email", Type: "string", Tags: `json:"email" field.name:"email" field.type:"varchar(255)"`},
			{Name: "gender", Type: "string", Tags: `json:"gender" field.name:"gender" field.type:"varchar(255)"`},
		},
	}
	output := structInfo.String()

	splitedoutput := strings.Split(output, "\n")
	var str []string

	for _, element := range splitedoutput {
		t := strings.TrimSpace(element)
		str = append(str, t)
	}

	joinedString := strings.Join(str, ", ")

	expected := `
Structure Name: Persons

		Name:	     Type:	Tags:
		-----	     -----	---------
		   id	     int64	json:"id" field.name:"id" field.type:"serial"
		 name	    string	json:"name" field.name:"name" field.type:"varchar(255)"
		email	    string	json:"email" field.name:"email" field.type:"varchar(255)"
	   gender	    string	json:"gender" field.name:"gender" field.type:"varchar(255)"
`

	splitedexpected := strings.Split(expected, "\n")
	var trimed []string

	for _, element := range splitedexpected {
		t := strings.TrimSpace(element)
		trimed = append(trimed, t)
	}

	expectedStr := strings.Join(trimed, ", ")

	if joinedString != expectedStr {
		t.Errorf(`Output not as expected. Output length = %d Expected length = %d
			--- # Output ---
			%s
			--- # Expected ---
			%s
			------------------
			`, len(output), len(expected), joinedString, expectedStr)
	}

}
