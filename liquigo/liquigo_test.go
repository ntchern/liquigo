package liquigo

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"gotest.tools/v3/assert"
)

func TestWithPostgres(t *testing.T) {

	db := postgresConn(t)
	defer db.Close()

	dropTable(db, "dbchangelog")
	dropTable(db, "dbchangelock")
	dropTable(db, "order")
	dropTable(db, "product")
	dropTable(db, "customer")
	dropTable(db, "category")
	dropTable(db, "address")

	t.Run("InitLogAndLockTables", func(t *testing.T) {
		assert.Assert(t, !isTableExist(t, db, "dbchangelog"))
		assert.Assert(t, !isTableExist(t, db, "dbchangelock"))
		err := Update(db, "../test-files/0-empty/_changelog.yaml")

		assert.NilError(t, err)
		assert.Assert(t, isTableExist(t, db, "dbchangelog"))
		assert.Assert(t, isTableExist(t, db, "dbchangelock"))

		// once more with tables already created
		err = Update(db, "../test-files/0-empty/_changelog.yaml")
		assert.NilError(t, err)
	})

	t.Run("RollbackTxIfFailedSQL", func(t *testing.T) {
		err := Update(db, "../test-files/1-invalid-sql/_changelog.yaml")
		assert.ErrorContains(t, err, "syntax error at or near \"BLAH\"")
		assert.Assert(t, !isTableExist(t, db, "order"))
	})

	t.Run("InitialUpdate", func(t *testing.T) {
		err := Update(db, "../test-files/2-initial-update/_changelog.yaml")
		assert.NilError(t, err)

		rows, err := findChangelogRecords(db)
		assert.NilError(t, err)
		assert.Equal(t, 2, len(rows))
		// first changeset
		assert.Equal(t, "order-product-table", rows[0].id)
		assert.Equal(t, "67fc014fff5646071881be4377e89095", rows[0].md5)
		assert.Equal(t, 1, rows[0].seqNo)
		// second changeset
		assert.Equal(t, "customer-table", rows[1].id)
		assert.Equal(t, "d7fa8981129ea6a2f7aedbf7b852095d", rows[1].md5)
		assert.Equal(t, 2, rows[1].seqNo)
	})

	t.Run("DoNothingOnRerun", func(t *testing.T) {
		err := Update(db, "../test-files/2-initial-update/_changelog.yaml")
		assert.NilError(t, err)

		rows, err := findChangelogRecords(db)
		assert.NilError(t, err)
		assert.Equal(t, 2, len(rows))
	})

	t.Run("DoNotAllowModifiedSQL", func(t *testing.T) {
		err := Update(db, "../test-files/3-modified-sql/_changelog.yaml")
		assert.Equal(t, err, ErrModifiedChangeset)
	})

	t.Run("ApplyAdditionalUpdate", func(t *testing.T) {
		err := Update(db, "../test-files/4-additional-update/_changelog.yaml")
		assert.NilError(t, err)

		rows, err := findChangelogRecords(db)
		assert.NilError(t, err)
		assert.Equal(t, 4, len(rows))
		// first changeset
		assert.Equal(t, "category-table", rows[2].id)
		assert.Equal(t, "ab9024d04bbef1296a26c1b6082ad45d", rows[2].md5)
		assert.Equal(t, 3, rows[2].seqNo)
		// second changeset
		assert.Equal(t, "address-table", rows[3].id)
		assert.Equal(t, "e0490e7d4ca84cc5c9b47e9525c9138d", rows[3].md5)
		assert.Equal(t, 4, rows[3].seqNo)
	})

}

func findChangelogRecords(db *sql.DB) ([]changeSetRecord, error) {
	result := []changeSetRecord{}
	rows, err := db.Query("SELECT id, md5, seq_no, tag FROM dbchangelog")
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var r changeSetRecord
		err = rows.Scan(&r.id, &r.md5, &r.seqNo, &r.tag)
		if err != nil {
			return result, err
		}
		result = append(result, r)
	}

	return result, nil
}
