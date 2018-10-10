package util

import (
	"fmt"
	"strconv"
	"strings"
)

//Postgres provides postgres compliant implementation of generated CRUD code
type Postgres struct {
	Entity
}

//SQLInsert returns SQL query for SQL Insert
func (s Postgres) SQLInsert() string {
	var (
		fields, placeholders []string
		index                = 1
	)

	for _, f := range s.Fields {
		fields = append(fields, f.schema.Field)
		placeholders = append(placeholders, "$"+strconv.Itoa(index))
		index++
	}

	for _, p := range s.Relationships {
		if p.Type == RelationshipTypeManyOne {
			fields = append(fields, fmt.Sprintf(`"%s"`, p.ThisID))
			placeholders = append(placeholders, "$"+strconv.Itoa(index))
			index++
		}
	}

	return fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) RETURNING id`,
		s.Table,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)
}

//BeforeInsert returns code executed before SQL statement is executed
func (s Postgres) BeforeInsert() []string {
	var before []string

	for _, f := range s.Fields {
		switch f.Name {
		case "CreatedAt", "UpdatedAt":
			before = append(before, fmt.Sprintf(`entity.%s = ptypes.TimestampNow()`, f.Name))
		}

		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		case "time":
			before = append(before, fmt.Sprintf(`%s, _ := pytpes.Timestamp(entity.%s)`, strings.ToLower(f.Name), f.Name))
		}

	}

	return before
}

//StructInsert returns list of fields to be used for insert statement
func (s Postgres) StructInsert() string {
	var (
		fields []string
	)

	for _, f := range s.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		default:
			fields = append(fields, fmt.Sprintf(`entity.%s`, f.Name))
		case "time":
			fields = append(fields, strings.ToLower(f.Name))
		}

	}
	return strings.Join(fields, ", ")
}

//AfterInsert returns code executed after SQL statement is executed
func (s Postgres) AfterInsert() []string {
	return []string{}
}
