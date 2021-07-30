package liquigo

import (
	"database/sql"
	"log"
	"os"
	"path"
)

func Update(db *sql.DB, file string) error {
	err := initTables(db)
	if err != nil {
		return err
	}

	err = acquireLock(db)
	if err != nil {
		return err
	}
	defer releaseLock(db)

	changelogFile, err := os.Open(file)
	if err != nil {
		return err
	}
	dir := path.Dir(file)
	files, err := parseChangelog(changelogFile)
	if err != nil {
		return err
	}
	var sets []changeset
	for _, f := range files.Files {
		changesetFile, err := os.Open(dir + "/" + f)
		if err != nil {
			return err
		}
		s, err := changesets(changesetFile)
		if err != nil {
			return err
		}
		sets = append(sets, s...)
	}

	for _, set := range sets {
		if err = set.apply(db, ""); err != nil {
			log.Println("Error applying changeset:", set.ID, err.Error())
			return err
		}
	}

	return nil
}

const (
	dbchangelog = `CREATE TABLE IF NOT EXISTS dbchangelog (
		id text PRIMARY KEY,
		applied_at timestamptz NOT NULL DEFAULT current_timestamp,
		seq_no serial,
		md5 text NOT NULL,
		tag text
	);`

	dbchangelock = `CREATE TABLE IF NOT EXISTS dbchangelock (
		id smallint PRIMARY KEY,
		locked_at timestamptz NOT NULL DEFAULT current_timestamp,
		locked_by text
	);`
	dbchangelockInsert = "INSERT INTO dbchangelock(id) VALUES (1)"
	dbchangelockDelete = "DELETE FROM dbchangelock"
)

func initTables(db *sql.DB) error {
	_, err := db.Exec(dbchangelog)
	if err != nil {
		return err
	}
	_, err = db.Exec(dbchangelock)
	if err != nil {
		return err
	}
	return nil
}

func acquireLock(db *sql.DB) error {
	// TODO with wait and retry?
	_, err := db.Exec(dbchangelockInsert)
	if err != nil {
		return err
	}
	return nil
}

func releaseLock(db *sql.DB) {
	db.Exec(dbchangelockDelete)
}
