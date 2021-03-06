package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/inflection"
)

const (
	// RelationshipTypeManyManyOwner represents a relationship of type Many(Owner)-Many
	RelationshipTypeManyManyOwner = "many-many-owner"

	// RelationshipTypeManyManyInverse represents a relationship of type Many-Many(owner)
	RelationshipTypeManyManyInverse = "many-many-inverse"

	// RelationshipTypeOneOne represents a relationship of type One-One
	RelationshipTypeOneOne = "one-one"

	RelationshipTypeOneOneOwner = "one-one-owner"

	// RelationshipTypeOneMany represents a relationship of type One-Many
	RelationshipTypeOneMany = "one-many"

	// RelationshipTypeManyOne represents a relationship of type Many-One
	RelationshipTypeManyOne = "many-one"
)

var (
	//ErrorRelationshipTargetEntityNotFound indicates that a field was not found during lookup request
	ErrorRelationshipTargetEntityNotFound = errors.New("relationship target entity not found")

	//ErrorInvalidRelationshipType indicates an invalid or unsupported relationship type
	ErrorInvalidRelationshipType = errors.New("relationship type is invalid")
)

// Relationship represents a relationship between this entity and another
type Relationship struct {
	// Entity is the name of the related entity
	Entity string `json:"entity"`

	// Type represents the type of this relationship
	Type string `json:"type"`

	// Name represents the property name to be used for this relationship
	Name string `json:"name"`

	// EditWidget represents widget information for this relationship
	EditWidget EditWidgetOpts `json:"edit_widget"`

	// JoinTable represents the other table in a many-many relationship
	JoinTable string `json:"-"`

	// ThisID represents the field in this entity used for this relationship (schema)
	ThisID string `json:"-"`

	// ThatID represents the field in the other entity used for this relationship (schema)
	ThatID string `json:"-"`

	RelSeeder RelationshipSeeder `json:"relation_seeder"`

	//related is a pointer to the related entity
	related *Entity
}

//RelationshipSeeder implements seeder option for relationship many-many
type RelationshipSeeder struct {
	MaxPerEntity int `json:"max"`
	FixPerEntity int `json:"fix"`
	MinPerEntity int `json:"min"`
}

func (p *Relationship) init(r *Recipe, e *Entity) {
	var isMany bool

	if rel, err := r.GetEntity(p.Entity); err == nil {
		p.related = rel
	} else {
		return
	}

	switch p.Type {
	case RelationshipTypeOneOne:
		p.JoinTable = ""
		p.ThisID = "id"
		p.ThatID = "id"

	case RelationshipTypeOneMany:
		isMany = true
		p.JoinTable = p.related.Table
		p.ThisID = "id"
		p.ThatID = inflection.Singular(e.Table) + "_id"
		if p.Name != "" {
			p.ThatID = strings.ToLower(p.Name) + "_" + p.ThatID
		}

	case RelationshipTypeManyOne:
		p.JoinTable = ""
		p.ThisID = inflection.Singular(p.related.Table) + "_id"
		if p.Name != "" {
			p.ThisID = strings.ToLower(p.Name) + "_" + p.ThisID
		}
		p.ThatID = "id"

	case RelationshipTypeManyManyOwner, RelationshipTypeManyManyInverse:
		isMany = true
		if strings.Compare(e.Table, p.related.Table) == -1 { //ascending order
			p.JoinTable = e.Table + "_" + p.related.Table
		} else {
			p.JoinTable = p.related.Table + "_" + e.Table
		}
		if p.Name != "" {
			p.JoinTable += "_" + strings.ToLower(p.Name)
		}
		p.ThisID = strings.ToLower(p.Entity) + "_id"
		p.ThatID = inflection.Singular(e.Table) + "_id"
	}

	if p.Name == "" {
		if isMany {
			p.Name = inflection.Plural(strings.Title(strings.ToLower(p.Entity)))
		} else {
			p.Name += strings.Title(p.Entity)
		}
	}
}

//Validate checks the relationship for errors
func (p *Relationship) Validate(r *Recipe) error {
	if _, err := r.GetEntity(p.Entity); err != nil {
		return ErrorRelationshipTargetEntityNotFound
	}

	switch p.Type {
	case RelationshipTypeManyManyOwner,
		RelationshipTypeManyManyInverse,
		RelationshipTypeOneOne,
		RelationshipTypeOneMany,
		RelationshipTypeManyOne:
		//ok
	default:
		return ErrorInvalidRelationshipType
	}

	return nil
}

//ProtoDefinitions returns proto definition for this relationship
func (p *Relationship) ProtoDefinitions(index *int) []string {
	var definitions []string

	switch p.Type {
	case RelationshipTypeManyManyOwner,
		RelationshipTypeManyManyInverse,
		RelationshipTypeOneMany:

		definitions = append(definitions, fmt.Sprintf(`repeated %s %s = %d;`, p.related.Name, p.Name, *index))
		*index++
	case RelationshipTypeOneOne:
		definitions = append(definitions, fmt.Sprintf(`%s %s = %d;`, p.related.Name, p.Name, *index))
		*index++
	case RelationshipTypeManyOne:
		definitions = append(definitions, fmt.Sprintf(`string %sID = %d;`, p.Name, *index))
		*index++
		definitions = append(definitions, fmt.Sprintf(`%s %s = %d;`, p.related.Name, p.Name, *index))
		*index++
	}

	return definitions
}

//GetRelatedID returns the related
func (p *Relationship) GetRelatedID() string {
	switch p.Type {
	case RelationshipTypeManyOne, RelationshipTypeManyManyOwner, RelationshipTypeManyManyInverse:
		return p.Name + "ID"
	}

	return "ID"
}

//GetRelated returns the related
func (p *Relationship) GetRelated() *Entity {
	return p.related
}
