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
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf(`"%s"`, p.ThisID))
			placeholders = append(placeholders, "$"+strconv.Itoa(index))
			index++
		}
	}

	return fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
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
	var fields []string

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

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf("entity.%sID", p.Name))
		}
	}

	return strings.Join(fields, ", ")
}

//AfterInsert returns code executed after SQL statement is executed
func (s Postgres) AfterInsert() []string {
	return []string{}
}

//SQLGet returns SQL query for SQL Get
func (s Postgres) SQLGet() string {
	var fields []string

	for _, f := range s.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf(`t."%s"`, p.ThisID))
		}
	}

	return fmt.Sprintf(
		`SELECT %s FROM %s t WHERE t."id" = $1 ORDER BY t."id" ASC`,
		strings.Join(fields, ", "),
		s.Table,
	)
}

//BeforeGet returns code executed before SQL statement is executed
func (s Postgres) BeforeGet() []string {
	var before []string

	for _, f := range s.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		case "time":
			before = append(before, fmt.Sprintf(`var %s time.Time`, strings.ToLower(f.Name)))
		}

	}

	return before
}

//StructGet returns list of fields to be used for get statement
func (s Postgres) StructGet() string {
	var fields []string

	for _, f := range s.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		default:
			fields = append(fields, fmt.Sprintf(`&entity.%s`, f.Name))
		case "time":
			fields = append(fields, fmt.Sprintf("&%s", strings.ToLower(f.Name)))
		}

	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf("&entity.%sID", p.Name))
		}
	}

	return strings.Join(fields, ", ")
}

//AfterGet returns code executed after SQL statement is executed
func (s Postgres) AfterGet() []string {
	var after []string

	for _, f := range s.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		case "time":
			after = append(after, fmt.Sprintf(`entity.%s, _ = ptypes.TimestampProto(%s)`, strings.ToLower(f.Name), f.Name))
		}

	}

	return after
}

//SQLList returns SQL query for SQL List
func (s Postgres) SQLList() string {
	var fields []string

	for _, f := range s.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf(`t."%s"`, p.ThisID))
		}
	}

	return fmt.Sprintf(
		`SELECT %s FROM %s t`,
		strings.Join(fields, ", "),
		s.Table,
	)
}

//BeforeList returns code executed before SQL statement is executed
func (s Postgres) BeforeList() []string {
	var before []string

	for _, f := range s.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		case "time":
			before = append(before, fmt.Sprintf(`var %s time.Time`, strings.ToLower(f.Name)))
		}

	}

	return before
}

//OrderList returns list of fields to be used for list statement
func (s Postgres) OrderList() string {
	var fields []string

	for _, f := range s.Fields {
		fields = append(fields, f.schema.Field)
	}

	return `"` + strings.Join(fields, `","`) + `"`
}

//StructList returns list of fields to be used for list statement
func (s Postgres) StructList() string {
	var fields []string

	for _, f := range s.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		default:
			fields = append(fields, fmt.Sprintf(`&entity.%s`, f.Name))
		case "time":
			fields = append(fields, fmt.Sprintf("&%s", strings.ToLower(f.Name)))
		}

	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf("&entity.%sID", p.Name))
		}
	}

	return strings.Join(fields, ", ")
}

//AfterList returns code executed after SQL statement is executed
func (s Postgres) AfterList() []string {
	var after []string

	for _, f := range s.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		switch f.Type {
		case "time":
			after = append(after, fmt.Sprintf(`entity.%s, _ = ptypes.TimestampProto(%s)`, strings.ToLower(f.Name), f.Name))
		}

	}

	return after
}

//SQLDeleteSingle returns SQL query for SQL Delete Single
func (s Postgres) SQLDeleteSingle() string {
	return fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1`,
		s.Table,
	)
}

//SQLDeleteMany returns SQL query for SQL Delete Many
func (s Postgres) SQLDeleteMany() string {
	return fmt.Sprintf(
		`DELETE FROM %s`,
		s.Table,
	)
}

//SQLDeleteManyJoin returns SQL query for SQL Delete Many Join
func (s Postgres) SQLDeleteManyJoin() []string {
	var query []string

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyManyOwner:
			query = append(query, fmt.Sprintf("DELETE FROM %s WHERE %s IN (SELECT id FROM %s", p.JoinTable, p.ThatID, s.Table))
		}
	}

	return query
}

//SQLUpdate returns SQL query for SQL Update
func (s Postgres) SQLUpdate() string {
	var (
		fields []string
		index  = 2
	)

	for _, f := range s.Fields {
		fields = append(fields, fmt.Sprintf(`"%s" = $%d`, f.schema.Field, index))
		index++
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne, RelationshipTypeOneOne:
			fields = append(fields, fmt.Sprintf(`"%s" = $%d`, p.ThisID, index))
			index++
		}
	}

	return fmt.Sprintf(
		`UPDATE %s SET %s WHERE id = $1`,
		s.Table,
		strings.Join(fields, ", "),
	)
}

//BeforeUpdate returns code executed before SQL statement is executed
func (s Postgres) BeforeUpdate() []string {
	var before []string

	for _, f := range s.Fields {
		switch f.Name {
		case "UpdatedAt":
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

//StructUpdate returns list of fields to be used for update statement
func (s Postgres) StructUpdate() string {
	var fields []string

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

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne, RelationshipTypeOneOne:
			fields = append(fields, fmt.Sprintf("entity.%sID", p.Name))
		}
	}

	return strings.Join(fields, ", ")
}

//AfterUpdate returns code executed after SQL statement is executed
func (s Postgres) AfterUpdate() []string {
	return []string{}
}

//SQLMerge returns SQL query for SQL Merge
func (s Postgres) SQLMerge() string {
	var (
		updates, inserts, placeholders []string
		index                          = 1
	)

	for _, f := range s.Fields {
		inserts = append(inserts, f.schema.Field)
		updates = append(updates, fmt.Sprintf(`"%s" = $%d`, f.schema.Field, index))
		placeholders = append(placeholders, "$"+strconv.Itoa(index))
		index++
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne, RelationshipTypeOneOne:
			inserts = append(inserts, fmt.Sprintf(`"%s"`, p.ThisID))
			updates = append(updates, fmt.Sprintf(`"%s" = $%d`, p.ThisID, index))
			placeholders = append(placeholders, "$"+strconv.Itoa(index))
			index++
		}
	}

	return fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) 
		ON CONFLICT (id) DO UPDATE SET %s`,
		s.Table,
		strings.Join(inserts, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "),
	)
}

//BeforeMerge returns code executed before SQL statement is executed
func (s Postgres) BeforeMerge() []string {
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

//StructMerge returns list of fields to be used for merge statement
func (s Postgres) StructMerge() string {
	var fields []string

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

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf("entity.%sID", p.Name))
		}
	}

	return strings.Join(fields, ", ")
}

//AfterMerge returns code executed after SQL statement is executed
func (s Postgres) AfterMerge() []string {
	return []string{}
}
