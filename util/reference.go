package util

// Reference represents a reference to another entity
type Reference struct {
	// IDField represents field information for referenced entity ID field
	IDField Field `json:"id_field"`

	// TypeField represents field information for referenced entity Type field
	TypeField Field `json:"type_field"`
}

func (c *Reference) init() {
	c.IDField.init()
	c.TypeField.init()
}

//Validate checks the reference definition for errors
func (c *Reference) Validate() error {
	if err := c.IDField.Validate(); err != nil {
		return err
	}

	if err := c.TypeField.Validate(); err != nil {
		return err
	}

	return nil
}
