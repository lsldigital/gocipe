package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

func generateLoadRelatedManyMany(entities map[string]util.Entity, entity util.Entity, rel util.Relationship) (string, error) {
	var sqlfields, structfields, before, after []string

	sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, "id"))
	structfields = append(structfields, fmt.Sprintf("&entity.%s", "ID"))
	related := entities[rel.Entity]

	thisType, _ := util.GetPrimaryKeyDataType(entity.PrimaryKey)
	thatType, _ := util.GetPrimaryKeyDataType(related.PrimaryKey)

	for _, field := range related.Fields {
		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("var %s time.Time", prop))
			structfields = append(structfields, fmt.Sprintf("&%s", prop))
			after = append(after, fmt.Sprintf("entity.%s, _ = ptypes.TimestampProto(%s)", field.Property.Name, prop))
		} else {
			structfields = append(structfields, fmt.Sprintf("&entity.%s", field.Property.Name))
		}
		sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, field.Schema.Field))
	}

	return util.ExecuteTemplate("crud/partials/loadrelated_manymany.go.tmpl", struct {
		ThisEntity   string
		Funcname     string
		ThatEntity   string
		SQLFields    string
		StructFields string
		JoinTable    string
		ThisType     string
		ThatTable    string
		ThatType     string
		ThisID       string
		ThatID       string
		PropertyName string
		Before       []string
		After        []string
		HasPreHook   bool
		HasPostHook  bool
	}{
		ThisEntity:   entity.Name,
		Funcname:     util.RelFuncName(rel),
		ThatEntity:   related.Name,
		JoinTable:    rel.JoinTable,
		ThisType:     thisType,
		ThatTable:    related.Table,
		ThatType:     thatType,
		ThisID:       strings.ToLower(entity.Name) + "_id",
		ThatID:       strings.ToLower(related.Name) + "_id",
		PropertyName: rel.Name,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		Before:       before,
		After:        after,
		HasPreHook:   entity.CrudHook.PreList,
		HasPostHook:  entity.CrudHook.PostList,
	})
}

func generateLoadRelatedOneMany(entities map[string]util.Entity, entity util.Entity, rel util.Relationship) (string, error) {
	var sqlfields, structfields, before, after []string

	sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, "id"))
	structfields = append(structfields, fmt.Sprintf("&entity.%s", "ID"))
	related := entities[rel.Entity]

	thisType, _ := util.GetPrimaryKeyDataType(entity.PrimaryKey)
	thatType, _ := util.GetPrimaryKeyDataType(related.PrimaryKey)

	for _, field := range related.Fields {
		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("var %s time.Time", prop))
			structfields = append(structfields, fmt.Sprintf("&%s", prop))
			after = append(after, fmt.Sprintf("entity.%s, _ = ptypes.TimestampProto(%s)", field.Property.Name, prop))
		} else {
			structfields = append(structfields, fmt.Sprintf("&entity.%s", field.Property.Name))
		}
		sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, field.Schema.Field))
	}

	for _, rel := range related.Relationships {
		if rel.Type == util.RelationshipTypeManyOne {
			sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, rel.ThisID))
			structfields = append(structfields, fmt.Sprintf("&entity.%sID", rel.Name))
		}
	}

	return util.ExecuteTemplate("crud/partials/loadrelated_onemany.go.tmpl", struct {
		ThisEntity   string
		Funcname     string
		ThatEntity   string
		SQLFields    string
		StructFields string
		ThisType     string
		ThatTable    string
		ThatType     string
		ThisID       string
		ThatID       string
		PropertyName string
		Before       []string
		After        []string
		HasPreHook   bool
		HasPostHook  bool
	}{
		ThisEntity:   entity.Name,
		Funcname:     util.RelFuncName(rel),
		ThatEntity:   related.Name,
		ThisType:     thisType,
		ThatTable:    related.Table,
		ThatType:     thatType,
		ThisID:       strings.ToLower(entity.Name) + "_id",
		ThatID:       "id",
		PropertyName: rel.Name,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		Before:       before,
		After:        after,
		HasPreHook:   entity.CrudHook.PreList,
		HasPostHook:  entity.CrudHook.PostList,
	})
}

func generateLoadRelatedManyOne(entities map[string]util.Entity, entity util.Entity, rel util.Relationship) (string, error) {
	var sqlfields, structfields, before, after []string

	sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, "id"))
	structfields = append(structfields, fmt.Sprintf("&thatEntity.%s", "ID"))
	related := entities[rel.Entity]

	thisType, _ := util.GetPrimaryKeyDataType(entity.PrimaryKey)
	thatType, _ := util.GetPrimaryKeyDataType(related.PrimaryKey)

	for _, field := range related.Fields {
		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("var %s time.Time", prop))
			structfields = append(structfields, fmt.Sprintf("&%s", prop))
			after = append(after, fmt.Sprintf("entity.%s, _ = ptypes.TimestampProto(%s)", field.Property.Name, prop))
		} else {
			structfields = append(structfields, fmt.Sprintf("&thatEntity.%s", field.Property.Name))
		}
		sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, field.Schema.Field))
	}

	return util.ExecuteTemplate("crud/partials/loadrelated_manyone.go.tmpl", struct {
		ThisEntity   string
		Funcname     string
		ThatEntity   string
		SQLFields    string
		StructFields string
		ThisID       string
		ThisType     string
		ThatTable    string
		ThatType     string
		PropertyName string
		Before       []string
		After        []string
		HasPreHook   bool
		HasPostHook  bool
	}{
		ThisEntity:   entity.Name,
		Funcname:     util.RelFuncName(rel),
		ThatEntity:   related.Name,
		ThisID:       rel.Name + "ID",
		ThisType:     thisType,
		ThatTable:    related.Table,
		ThatType:     thatType,
		PropertyName: rel.Name,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		Before:       before,
		After:        after,
		HasPreHook:   entity.CrudHook.PreList,
		HasPostHook:  entity.CrudHook.PostList,
	})
}
