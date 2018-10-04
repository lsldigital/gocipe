package util

import "errors"

var (
	//ErrorReferenceNameEmpty indicates that the reference has no name
	ErrorReferenceNameEmpty = errors.New("reference name is empty")
)

// Reference represents a reference to another entity
type Reference struct {
	//Name is the reference field name
	Name string

	// iDField represents field information for referenced entity ID field
	idField Field

	// typeField represents field information for referenced entity Type field
	typeField Field
}

func (c *Reference) init() {
	c.idField, c.typeField = fieldReferenceMakeFields(c.Name)
	c.idField.init()
	c.typeField.init()
}

//Validate checks the reference definition for errors
func (c *Reference) Validate() error {
	if c.Name == "" {
		return ErrorReferenceNameEmpty
	}

	if err := c.idField.Validate(); err != nil {
		return err
	}

	if err := c.typeField.Validate(); err != nil {
		return err
	}

	return nil
}
