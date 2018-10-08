package util

import (
	"fmt"
	"strconv"
	"strings"
)

type Postgres struct {
	Entity
}

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
			before = append(before, fmt.Sprintf(`%s, + := pytpes.Timestamp(entity.%s)`, strings.ToLower(f.Name), f.Name))
		}

	}

	return before
}

func (s Postgres) StructInsert() string {
	var (
		fields []string
	)

	for _, f := range s.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		default:
			fields = append(fields, fmt.Sprintf(`"%s"`, f.schema.Field))
		case "time":
			fields = append(fields, fmt.Sprintf(`"%s"`, strings.ToLower(f.Name)))
		}

	}
	return ""
}

func (s Postgres) AfterInsert() string {
	return ""
}
