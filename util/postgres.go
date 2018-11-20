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

	for _, c := range s.References {
		fields = append(fields, c.IDField.schema.Field)
		placeholders = append(placeholders, "$"+strconv.Itoa(index))
		index++

		fields = append(fields, c.TypeField.schema.Field)
		placeholders = append(placeholders, "$"+strconv.Itoa(index))
		index++
	}

	return fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		s.Table,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)
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

	for _, c := range s.References {
		fields = append(fields, fmt.Sprintf(`t."%s"`, c.IDField.schema.Field), fmt.Sprintf(`t."%s"`, c.TypeField.schema.Field))
	}

	return fmt.Sprintf(
		`SELECT %s FROM %s t WHERE t."id" = $1 ORDER BY t."id" ASC`,
		strings.Join(fields, ", "),
		s.Table,
	)
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

	for _, c := range s.References {
		fields = append(fields, fmt.Sprintf(`t."%s"`, c.IDField.schema.Field), fmt.Sprintf(`t."%s"`, c.TypeField.schema.Field))
	}

	return fmt.Sprintf(
		`SELECT %s FROM %s t`,
		strings.Join(fields, ", "),
		s.Table,
	)
}

//OrderList returns list of fields to be used for list statement
func (s Postgres) OrderList() string {
	var fields []string

	for _, f := range s.Fields {
		fields = append(fields, f.schema.Field)
	}

	return `"` + strings.Join(fields, `","`) + `"`
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

	for _, c := range s.References {
		fields = append(fields, fmt.Sprintf(`"%s" = $%d`, c.IDField.schema.Field, index))
		index++

		fields = append(fields, fmt.Sprintf(`"%s" = $%d`, c.TypeField.schema.Field, index))
		index++
	}

	return fmt.Sprintf(
		`UPDATE %s SET %s WHERE id = $1`,
		s.Table,
		strings.Join(fields, ", "),
	)
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

	for _, c := range s.References {
		inserts = append(inserts, c.IDField.schema.Field)
		updates = append(updates, fmt.Sprintf(`"%s" = $%d`, c.IDField.schema.Field, index))
		placeholders = append(placeholders, "$"+strconv.Itoa(index))
		index++

		inserts = append(inserts, c.TypeField.schema.Field)
		updates = append(updates, fmt.Sprintf(`"%s" = $%d`, c.TypeField.schema.Field, index))
		placeholders = append(placeholders, "$"+strconv.Itoa(index))
		index++
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

//SQLLoadManyMany returns SQL query for SQL Load many many
func (s Postgres) SQLLoadManyMany(rel Relationship) string {
	var fields []string

	related := rel.related

	for _, f := range related.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	// many-one fields when loading entities in many-many-inverse relationships
	for _, r := range related.Relationships {
		if r.Type == RelationshipTypeManyOne {
			fields = append(fields, fmt.Sprintf(`t."%s"`, r.ThisID))
		}
	}

	return fmt.Sprintf(
		`SELECT j.%s, %s FROM %s t 
		INNER JOIN %s j ON t.id = j.%s
		WHERE j.%s IN`,
		rel.ThatID,
		strings.Join(fields, ", "),
		related.Table,
		rel.JoinTable,
		rel.ThisID,
		rel.ThatID,
	)
}

//SQLLoadManyOne returns SQL query for SQL Load many one
func (s Postgres) SQLLoadManyOne(rel Relationship) string {
	var fields []string

	related := rel.related

	for _, f := range related.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	for _, rel := range related.Relationships {
		if rel.Type == RelationshipTypeManyOne {
			fields = append(fields, fmt.Sprintf(`t."%s"`, rel.ThisID))
		}
	}

	for _, c := range related.References {
		fields = append(fields,
			fmt.Sprintf(`t."%s"`, c.IDField.schema.Field),
			fmt.Sprintf(`t."%s"`, c.TypeField.schema.Field))
	}

	return fmt.Sprintf(
		`SELECT %s FROM %s t WHERE t."%s" IN`,
		strings.Join(fields, ", "),
		related.Table,
		rel.ThatID,
	)
}

//SQLLoadOneMany returns SQL query for SQL Load one many
func (s Postgres) SQLLoadOneMany(rel Relationship) string {
	var fields []string

	related := rel.related

	for _, f := range related.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	for _, rel := range related.Relationships {
		if rel.Type == RelationshipTypeManyOne {
			fields = append(fields, fmt.Sprintf(`t."%s"`, rel.ThisID))
		}
	}

	for _, c := range related.References {
		fields = append(fields,
			fmt.Sprintf(`t."%s"`, c.IDField.schema.Field),
			fmt.Sprintf(`t."%s"`, c.TypeField.schema.Field))
	}

	return fmt.Sprintf(
		`SELECT %s FROM %s t WHERE t."%s" IN`,
		strings.Join(fields, ", "),
		related.Table,
		rel.ThatID,
	)
}

//SQLSaveManyManyOwnerDelete returns SQL query for SQL Save many many owner
func (s Postgres) SQLSaveManyManyOwnerDelete(rel Relationship) string {
	return fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1`,
		rel.JoinTable,
		rel.ThatID,
	)
}

//SQLSaveManyManyOwnerInsert returns SQL query for SQL Save many many owner
func (s Postgres) SQLSaveManyManyOwnerInsert(rel Relationship) string {
	return fmt.Sprintf(
		`INSERT INTO %s (%s, %s) VALUES ($1, $2)`,
		rel.JoinTable,
		rel.ThatID,
		rel.ThisID,
	)
}
