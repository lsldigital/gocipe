package util

import "errors"

var (
	//ErrorReferenceNameEmpty indicates that the reference has no name
	ErrorReferenceNameEmpty = errors.New("reference name is empty")
)

// Reference represents a reference to another entity
type Reference struct {
	//Name is the reference field name
	Name string `json:"name"`

	// IDField represents field information for referenced entity ID field
	IDField Field `json:"-"`

	// TypeField represents field information for referenced entity Type field
	TypeField Field `json:"-"`
}

func (c *Reference) init() {
	c.IDField, c.TypeField = fieldReferenceMakeFields(c.Name)
	c.IDField.init()
	c.TypeField.init()
}

//Validate checks the reference definition for errors
func (c *Reference) Validate() error {
	if c.Name == "" {
		return ErrorReferenceNameEmpty
	}

	if err := c.IDField.Validate(); err != nil {
		return err
	}

	if err := c.TypeField.Validate(); err != nil {
		return err
	}

	return nil
}
