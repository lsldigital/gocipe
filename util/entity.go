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
	// ActionAdd for action Add
	ActionAdd = "add"
	// ActionRemove for action Remove
	ActionRemove = "remove"
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

	// CreatedAt indicates whether to add a predefined CreatedAt timestamp field in the entity
	CreatedAt bool `json:"created_at"`

	// UpdatedAt indicates whether to add a predefined CreatedAt timestamp field in the entity
	UpdatedAt bool `json:"updated_at"`

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

	// ContentBuilder represents the Lardwaz module for content building
	ContentBuilder ContentBuilderOpts `json:"content_builder"`

	// Table is the name of the database table for the entity
	Table string `json:"-"`

	SeedCount int `json:"seed_count"`

	// fields is a map for random access of fields contained in entity
	fields map[string]*Field
}

func (e *Entity) init(r *Recipe) {
	var defaultFields = []Field{fieldID, fieldStatus}
	e.Table = inflection.Plural(strings.ToLower(e.Name))

	if e.Slug != "" {
		defaultFields = append(defaultFields, fieldSlug)
	}

	if e.ContentBuilder.Generate {
		defaultFields = append(defaultFields, contentField)
	}

	if r.Admin.Auth.Generate {
		defaultFields = append(defaultFields, fieldUserID)
	}

	if e.CreatedAt {
		defaultFields = append(defaultFields, fieldCreatedAt)
	}

	if e.UpdatedAt {
		defaultFields = append(defaultFields, fieldUpdatedAt)
	}

	e.Fields = append(defaultFields, e.Fields...)

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

	for i := range e.References {
		c := &e.References[i]

		c.init()

		// Fill in reference type field edit widget options for Card entity
		if r.Decks.Generate && e.Name == "Card" {
			// TODO: Options based on current deck selected (@runtime) ?
			for _, t := range r.Entities {
				if t.Name == "Card" || t.Name == "CardSchedule" {
					continue
				}
				c.TypeField.EditWidget.Options = append(c.TypeField.EditWidget.Options,
					EditWidgetOption{Text: t.Name, Value: inflection.Plural(t.Name)},
				)
			}
		}
	}

	e.fields = make(map[string]*Field)
	for i := range e.Fields {
		f := &e.Fields[i]

		// Fill in deckmachinename field edit widget options for Card entity
		if r.Decks.Generate && e.Name == "Card" && f.Name == "DeckMachineName" {
			for _, d := range r.Decks.Decks {
				f.EditWidget.Options = append(f.EditWidget.Options,
					EditWidgetOption{Text: d.Label, Value: d.Name},
				)
			}
		}

		f.init()
		e.fields[f.Name] = f
	}

	if e.LabelField == "" {
		for _, p := range preferredLabelFields { //check if preferred label field present
			if f, err := e.GetField(p); err == nil && f.Type == fieldTypeStr {
				e.LabelField = f.schema.Field
				break
			}
		}

		if e.LabelField == "" { //find first string field for label
			for _, f := range e.Fields {
				if f.Type == fieldTypeStr {
					e.LabelField = f.schema.Field
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

	for _, c := range e.References {
		filtersString = append(filtersString, c.IDField.schema.Field, c.TypeField.schema.Field)
		filters.HasString = true
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

//GetForeignKeyFields returns definition of related foreign key fields (used by schema)
func (e *Entity) GetForeignKeyFields() []string {
	var related []string

	for _, p := range e.Relationships {
		if p.Type == RelationshipTypeManyOne {
			related = append(related, fmt.Sprintf(`"%s" TEXT NOT NULL`, p.ThisID)) //TODO SQL-dialect sensitive
		}
	}

	return related
}

//GetReferenceFields returns definition of reference fields (used by schema)
func (e *Entity) GetReferenceFields() []string {
	var referenced []string

	for _, p := range e.References {
		referenced = append(referenced,
			fmt.Sprintf(`"%s" TEXT NOT NULL`, p.IDField.schema.Field),
			fmt.Sprintf(`"%s" TEXT NOT NULL`, p.TypeField.schema.Field)) //TODO SQL-dialect sensitive
	}

	return referenced
}

//GetRelatedTables returns definition of related tables; typically due to many-many relationships (used by schema)
func (e *Entity) GetRelatedTables() []Relationship {
	var related []Relationship
	for _, p := range e.Relationships {
		switch p.Type {
		default:
			continue
		case RelationshipTypeManyManyInverse, RelationshipTypeManyManyOwner:
			if p.related == nil || strings.Compare(e.Name, p.related.Name) < 0 {
				continue
			}
		}

		related = append(related, p)
	}

	return related
}

//GetLabelFieldName returns field name for label field
func (e *Entity) GetLabelFieldName() string {
	for _, f := range e.Fields {
		if f.schema.Field == e.LabelField {
			return f.Name
		}
	}

	return ""
}

//HasCrudHooks returns true if any of crud hooks is enabled
func (e *Entity) HasCrudHooks() bool {
	switch true {
	case e.CrudHooks.PreSave,
		e.CrudHooks.PostSave,
		e.CrudHooks.PreRead,
		e.CrudHooks.PostRead,
		e.CrudHooks.PreList,
		e.CrudHooks.PostList,
		e.CrudHooks.PreDeleteSingle,
		e.CrudHooks.PostDeleteSingle,
		e.CrudHooks.PreDeleteMany,
		e.CrudHooks.PostDeleteMany:
		return true

	}
	return false
}

//HasAdminHooks returns true if any of admin hooks is enabled
func (e *Entity) HasAdminHooks() bool {
	switch true {
	case e.Admin.Hooks.PreCreate,
		e.Admin.Hooks.PostCreate,
		e.Admin.Hooks.PreRead,
		e.Admin.Hooks.PostRead,
		e.Admin.Hooks.PreList,
		e.Admin.Hooks.PostList,
		e.Admin.Hooks.PreUpdate,
		e.Admin.Hooks.PostUpdate,
		e.Admin.Hooks.PreDelete:
		return true

	}
	return false
}

// HasTimestamp returns true if the entity has at least 1 field of type timestamp
func (e *Entity) HasTimestamp() bool {
	for _, f := range e.Fields {
		if f.Type == fieldTypeTime {
			return true
		}
	}
	return false
}

//GetProtoFields returns list of protobuf field definitions for this entity
func (e *Entity) GetProtoFields() []string {
	var (
		fields []string
		index  = 1
	)

	for _, f := range e.Fields {
		fields = append(fields, f.ProtoDefinition(&index))
	}

	for _, p := range e.Relationships {
		fields = append(fields, p.ProtoDefinitions(&index)...)
	}

	for _, c := range e.References {
		fields = append(fields, c.IDField.ProtoDefinition(&index), c.TypeField.ProtoDefinition(&index))
	}

	return fields
}

//GetStruct returns list of fields to be used for 'op' statement
func (e *Entity) GetStruct(op string) string {
	var fields []string
	switch op {
	case "get", "list":
		//
	case "update":
		fields = append(fields, "entity.ID")
	}

	for _, f := range e.Fields {
		//some fields require preprocessing
		//they will be assigned to a variable, use that instead of the property name
		if f.Type == fieldTypeTime {
			switch op {
			case "get", "list":
				fields = append(fields, fmt.Sprintf("&%s", strings.ToLower(f.Name)))
			case "insert", "merge", "update":
				fields = append(fields, strings.ToLower(f.Name))
			}
		} else {
			switch op {
			case "get", "list":
				fields = append(fields, fmt.Sprintf(`&entity.%s`, f.Name))
			case "insert", "merge", "update":
				fields = append(fields, fmt.Sprintf(`entity.%s`, f.Name))
			}
		}
	}

	for _, p := range e.Relationships {
		if p.Type == RelationshipTypeManyOne {
			switch op {
			case "get", "list":
				fields = append(fields, fmt.Sprintf("&entity.%sID", p.Name))
			case "insert", "merge", "update":
				fields = append(fields, fmt.Sprintf("entity.%sID", p.Name))
			}
		}
	}

	for _, c := range e.References {
		switch op {
		case "get", "list":
			fields = append(fields, fmt.Sprintf("&entity.%s", c.IDField.Name), fmt.Sprintf("&entity.%s", c.TypeField.Name))
		case "insert", "merge", "update":
			fields = append(fields, fmt.Sprintf("entity.%s", c.IDField.Name), fmt.Sprintf("entity.%s", c.TypeField.Name))
		}
	}

	return strings.Join(fields, ", ")
}

//GetFileFields returns list of file fields
func (e *Entity) GetFileFields() []FileField {
	var fields []FileField

	if e.ContentBuilder.Generate && e.ContentBuilder.AllowUpload {
		for _, u := range e.ContentBuilder.UploadTypes {
			uploadType := strings.Title(u)
			fields = append(fields, FileField{
				ConfigName:     e.Name + "Content" + uploadType + "UploadOpts",
				Destination:    strings.ToLower(e.Name + "/content/" + uploadType),
				EntityName:     e.Name,
				FieldName:      uploadType,
				ContentBuilder: true,
			})
		}
	}

	for _, f := range e.Fields {
		switch f.EditWidget.Type {
		case WidgetTypeFile, WidgetTypeImage:
			fields = append(fields, FileField{
				ConfigName:     e.Name + f.Name + "UploadOpts",
				Destination:    strings.ToLower(e.Name + "/" + f.Name),
				EntityName:     e.Name,
				FieldName:      f.Name,
				SchemaName:     f.schema.Field,
				ContentBuilder: false,
			})
		}
	}

	return fields
}

//HasFileFields returns whether entity has file fields
func (e *Entity) HasFileFields() bool {
	if e.ContentBuilder.AllowUpload && len(e.ContentBuilder.UploadTypes) != 0 {
		return true
	}

	for _, f := range e.Fields {
		switch f.EditWidget.Type {
		case WidgetTypeFile, WidgetTypeImage:
			return true
		}
	}

	return false
}
