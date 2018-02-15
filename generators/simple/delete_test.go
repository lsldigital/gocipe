package simple

import (
	"strings"
	"testing"
)

func TestGenerateDelete(t *testing.T) {
	name := "Person"

	output, err := GenerateDelete(name)
	expected := `
//Delete delete single Person entity from database
func Delete(db *sql.DB, id int) error {
	stmt, err := db.Prepare("DELETE FROM ` + "`persons`" + ` WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
`

	if err != nil {
		t.Errorf("Got error: %s", err)
	} else if strings.Compare(output, expected) != 0 {
		t.Errorf(`Output not as expected. Output length = %d Expected length = %d
--- # Output ---
%s
--- # Expected ---
%s
------------------
`, len(output), len(expected), output, expected)
	}
}
