package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/inflection"
)

const (
	// StatusDraft "D" for status Draft
	StatusDraft = "draft"
	// StatusSaved "S" for status Saved
	StatusSaved = "saved"
	// StatusUnpublished "U" for status Unpublished
	StatusUnpublished = "unpublished"
	// StatusPublished "P" for status Published
	StatusPublished = "published"
)

var (
	preferredLabelFields = []string{"Name", "Title", "Description", "Summary", "BannerType"}

	//ErrorEntityNameBlank indicates an entity's name has not been defined (in gocipe.json)
	ErrorEntityNameBlank = errors.New("entity does not have a name")

	//ErrorEntityInvalidPrimaryKey indicates an invalid priamry key type
	ErrorEntityInvalidPrimaryKey = errors.New("entity primary key type is not supported")

	//ErrorEntityFieldNotFound indicates that a field was not found during lookup request
	ErrorEntityFieldNotFound = errors.New("entity does not contain requested field")

	//ErrorEntitySlugNotString indicates that the slug field for an entity is not of string as required
	ErrorEntitySlugNotString = errors.New("entity slug field must be a string")

	//ErrorEntityLabelFieldNotString indicates that the slug field for an entity is not of string as required
	ErrorEntityLabelFieldNotString = errors.New("entity label field must be a string")
)

// AdminFilters used for List
type AdminFilters struct {
	HasBool, HasString, HasDate             bool
	BoolFilters, StringFilters, DateFilters string
}

// Entity represents a single entity to be generated
type Entity struct {
	// Name is the name of the entity
	Name string `json:"name"`

	// PrimaryKey indicates the nature of the primary key: serial (auto incremented number), uuid (auto generated string), int or string
	PrimaryKey string `json:"primary_key"`

	// TableConstraints represents an array of table constraints for the table definition
	TableConstraints []string `json:"table_constraints"`

	// Slug indicates the field to be used for entity slug
	Slug string `json:"slug"`

	// Description is a description of the entity
	Description string `json:"description"`

	// Fields is a list of fields for the entity
	Fields []Field `json:"fields"`

	// LabelField indicates which field is used as label for the entity
	LabelField string `json:"label_field"`

	// Relationships represents relationship information between this entity and others
	Relationships []Relationship `json:"relationships"`

	// References represents reference information for other entities
	References []Reference `json:"references"`

	// Schema describes options for Schema generation - overrides recipe level Schema config
	Schema *SchemaOpts `json:"schema"`

	// Crud describes options for CRUD generation - overrides recipe level Crud config
	CrudHooks *CrudHooks `json:"crud"`

	// Admin describes options for Admin generation - overrides recipe level Admin config
	Admin *AdminOpts `json:"admin"`

	// Vuetify describes options for Vuetify generation - overrides recipe level Vuetify config
	Vuetify VuetifyEntityOpts `json:"vuetify"`

	// DefaultSort is a sort string used while generating List() method in CRUD
	DefaultSort string `json:"default_sort"`

	// Table is the name of the database table for the entity
	Table string `json:"-"`

	//fields is a map for random access of fields contained in entity
	fields map[string]*Field
}

func (e *Entity) init(r *Recipe) {
	e.Table = inflection.Plural(strings.ToLower(e.Name))

	e.Fields = append(e.Fields, fieldStatus)

	if r.Admin.Auth.Generate {
		e.Fields = append(e.Fields, fieldUserID)
	}

	if e.CrudHooks == nil {
		e.CrudHooks = &r.Crud.Hooks
	}

	if e.Admin == nil {
		e.Admin = &r.Admin
	}

	if e.PrimaryKey == "" {
		e.PrimaryKey = PrimaryKeySerial
	}

	if e.Vuetify.Icon == "" {
		e.Vuetify.Icon = "dashboard"
	}

	if e.DefaultSort == "" {
		e.DefaultSort = `t."id" DESC`
	}

	for i := range e.Fields {
		f := &e.Fields[i]
		f.init()
		e.fields[f.Name] = f
	}

	if e.Slug != "" {
		e.Fields = append(e.Fields, fieldSlug)
	}

	if e.LabelField == "" {
		for _, p := range preferredLabelFields { //check if preferred label field present
			if f, err := e.GetField(p); err == nil && f.Type == fieldTypeStr {
				e.LabelField = f.Name
				break
			}
		}

		if e.LabelField == "" { //find first string field for label
			for _, f := range e.Fields {
				if f.Type == fieldTypeStr {
					e.LabelField = f.Name
					break
				}
			}
		}
	}
}

//Validate checks the entity for errors
func (e *Entity) Validate() error {
	if e.Name == "" {
		return ErrorEntityNameBlank
	}

	switch e.PrimaryKey {
	case PrimaryKeyString, PrimaryKeyUUID:
		//all ok
	default:
		return ErrorEntityInvalidPrimaryKey
	}

	if e.LabelField == "" {
	} else if f, err := e.GetField(e.LabelField); err != nil {
		return err
	} else if f.Type != fieldTypeStr {
		return ErrorEntitySlugNotString
	}

	if e.Slug == "" {
	} else if f, err := e.GetField(e.Slug); err != nil {
		return err
	} else if f.Type != fieldTypeStr {
		return ErrorEntitySlugNotString
	}

	return nil
}

//GetField returns a field by name from the entity
func (e *Entity) GetField(name string) (*Field, error) {
	if f, ok := e.fields[name]; ok {
		return f, nil
	}
	return nil, ErrorEntityFieldNotFound
}

//GetAdminFilters returns a list of filters applicable to this entity (service admin)
func (e *Entity) GetAdminFilters() AdminFilters {
	var (
		filters                                 AdminFilters
		filtersBool, filtersString, filtersDate []string
	)

	for _, field := range e.Fields {
		switch field.Type {
		case fieldTypeBool:
			filtersBool = append(filtersBool, field.schema.Field)
			filters.HasBool = true
		case fieldTypeStr:
			filtersString = append(filtersString, field.schema.Field)
			filters.HasString = true
		case fieldTypeTime:
			filtersDate = append(filtersDate, field.schema.Field)
			filters.HasDate = true
		}
	}

	for _, rel := range e.Relationships {
		switch rel.Type {
		case RelationshipTypeOneOne, RelationshipTypeManyOne:
			filtersString = append(filtersString, rel.ThisID)
			filters.HasString = true
		}
	}

	if len(filtersBool) != 0 {
		filters.BoolFilters = `"` + strings.Join(filtersBool, `","`) + `"`
	}

	if len(filtersString) != 0 {
		filters.StringFilters = `"` + strings.Join(filtersString, `","`) + `"`
	}

	if len(filtersDate) != 0 {
		filters.DateFilters = `"` + strings.Join(filtersDate, `","`) + `"`
	}

	return filters
}

//GetFileFieldsDefinition returns file field definitions (used by admin service)
func (e *Entity) GetFileFieldsDefinition() []string {
	var fileFields []string
	for _, field := range e.Fields {
		switch field.EditWidget.Type {
		case WidgetTypeFile, WidgetTypeImage:
			tpl := strings.Join([]string{
				`case "%s":`,
				"options = &%s%sUploadOpts",
				`fieldname = "%s"`,
			}, "\n")
			fileFields = append(fileFields, fmt.Sprintf(tpl, field.Name, e.Name, field.Name, field.schema.Field))
		}
	}

	return fileFields
}
