package liquigo

import (
	"bufio"
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	commentChangeset = "-- changeset"
	commentMD5       = "-- md5"
	commentSplit     = "-- splitStatements"
	commentPrefix    = "-- "
)

var ErrModifiedChangeset = errors.New("liquigo: md5 does not match")

// Changeset represents a Changeset
type changeset struct {
	// ID of the changeset.
	ID string

	// MD5
	MD5 string

	// splitStatements
	splitStatements bool

	// Body of the changeset
	SQLs []string
}

// Changesets reads an input reader line by line, parsing it.
func changesets(r io.Reader) ([]changeset, error) {
	result := []changeset{}
	scanner := bufio.NewScanner(r)

	var set changeset
	var sb strings.Builder

	for scanner.Scan() {
		s := strings.Join(strings.Fields(scanner.Text()), " ")

		if strings.HasPrefix(s, commentChangeset) {
			if set.ID != "" {
				result = append(result, set)
			}
			set = changeset{
				ID:              parseID(s),
				splitStatements: true,
			}
			continue

		} else if strings.HasPrefix(s, commentMD5) {
			set.MD5 = parseMD5(s)
			continue

		} else if strings.HasPrefix(s, commentSplit) {
			b, err := parseSplit(s)
			if err != nil {
				return result, err
			}
			set.splitStatements = b
			continue

		} else if strings.HasPrefix(s, commentPrefix) {
			continue

		}

		s = strings.TrimSpace(strings.Split(s, commentPrefix)[0])
		if len(s) > 0 {
			if sb.Len() > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(s)
			if set.splitStatements && strings.HasSuffix(s, ";") {
				appendSQL(&set, sb)
				sb.Reset()
			}
		}
	}

	if sb.Len() > 0 {
		appendSQL(&set, sb)
	}
	if set.ID != "" {
		result = append(result, set)
	}
	return result, scanner.Err()
}

func appendSQL(set *changeset, sb strings.Builder) {
	set.SQLs = append(set.SQLs, sb.String())
	if set.MD5 == "" {
		m := md5.Sum([]byte(strings.Join(set.SQLs, " ")))
		set.MD5 = fmt.Sprintf("%x", m)
	}
}

type changeSetRecord struct {
	id        string
	appliedAt time.Time
	seqNo     int
	md5       string
	tag       string
}

func (set changeset) apply(db *sql.DB, tag string) error {
	toApply, err := toBeApplied(db, set.ID, set.MD5)
	if err != nil {
		return err
	}
	if !toApply {
		return nil
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, sql := range set.SQLs {
		_, err := tx.Exec(sql)
		if err != nil {
			return err
		}
	}

	dbRecord := changeSetRecord{
		id:        set.ID,
		appliedAt: time.Now(),
		md5:       set.MD5,
		tag:       tag,
	}

	_, err = tx.Exec(
		"INSERT INTO dbchangelog(id, applied_at, md5, tag) VALUES ($1, $2, $3, $4)",
		dbRecord.id, dbRecord.appliedAt, dbRecord.md5, dbRecord.tag,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Changeset applied: %s", dbRecord.id))
	return nil
}

func toBeApplied(db *sql.DB, id string, md5 string) (bool, error) {
	row := db.QueryRow("SELECT md5 FROM dbchangelog WHERE id = $1", id)
	found := changeSetRecord{}
	err := row.Scan(&found.md5)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	if found.md5 != md5 {
		return false, ErrModifiedChangeset
	}
	return false, nil
}

func parseID(input string) string {
	return strings.Trim(strings.TrimPrefix(input, commentChangeset), " \t")
}

func parseMD5(input string) string {
	return strings.Trim(strings.TrimPrefix(input, commentMD5), " \t")
}

func parseSplit(input string) (bool, error) {
	s := strings.Trim(strings.TrimPrefix(input, commentSplit), " \t")
	return strconv.ParseBool(s)
}
