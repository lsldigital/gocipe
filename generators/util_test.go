package generators

import (
	"strings"
	"testing"
)

func TestGetAbsPath(t *testing.T) {
	output, err := GetAbsPath("generators")

	expected := "/home/jeshta/go/src/lab.lsl.digital/gocipe/generators/generators"

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
