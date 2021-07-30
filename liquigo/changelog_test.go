package liquigo

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseChangelog(t *testing.T) {
	input := `
databaseChangeLog:
  - first
  - second
  - third
`

	changelog, err := parseChangelog(strings.NewReader(input))
	if err != nil {
		t.Fatalf("error parsing blocks: %v", err)
	}

	assert.Equal(t, 3, len(changelog.Files))
	assert.Equal(t, "first", changelog.Files[0])
}
