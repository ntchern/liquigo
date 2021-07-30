package liquigo

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestParse(t *testing.T) {
	input := `
	-- changeset a
	CREATE TABLE
		order; -- a comment

	CREATE TABLE
		product;

	-- changeset b
	-- md5 12345
	ALTER TABLE
		order;
	`

	sets, err := changesets(strings.NewReader(input))
	if err != nil {
		t.Fatalf("error parsing blocks: %v", err)
	}

	assert.Equal(t, 2, len(sets))
	assert.Equal(t, 2, len(sets[0].SQLs))
	assert.Equal(t, 1, len(sets[1].SQLs))

	assert.Equal(t, "a", sets[0].ID)
	assert.Equal(t, "CREATE TABLE order;", sets[0].SQLs[0])
	assert.Equal(t, "CREATE TABLE product;", sets[0].SQLs[1])
	assert.Equal(t, "5947e68449e2fd7cfbe20386e64023c9", sets[0].MD5)

	assert.Equal(t, "b", sets[1].ID)
	assert.Equal(t, "ALTER TABLE order;", sets[1].SQLs[0])
	assert.Equal(t, "12345", sets[1].MD5)
}
