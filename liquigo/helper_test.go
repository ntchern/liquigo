package liquigo

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const postgresVersion = "13.4-alpine"
const postgresPort = "5454"

func startPostgresContainer() error {
	return exec.Command("docker", "run", "-d", "--rm",
		"-e", "POSTGRES_PASSWORD=postgres",
		"--name", "postgres_liquigo",
		"-p", postgresPort+":5432",
		"postgres:"+postgresVersion).Start()
}

func postgresConn(t *testing.T) *sql.DB {
	startPostgresContainer()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s search_path=%s sslmode=disable",
		"localhost", postgresPort, "postgres", "postgres", "postgres", "public")
	deadline := time.Now().Add(10 * time.Second)
	var db *sql.DB
	for {
		db1, err := sql.Open("postgres", psqlInfo)
		if err == nil {
			db = db1
			break
		}
		if time.Now().After(deadline) {
			t.Fatal(err.Error())
			return nil
		}
	}
	for {
		r, err := db.Query("SELECT 1")
		if err == nil {
			r.Close()
			return db
		}
		if time.Now().After(deadline) {
			db.Close()
			t.Fatal(err.Error())
			return nil
		}
	}
}

func dropTable(db *sql.DB, table string) {
	db.Exec(fmt.Sprintf("DROP table \"%s\"", table))
}

func isTableExist(t *testing.T, db *sql.DB, table string) bool {
	sql := fmt.Sprintf("SELECT 1 FROM \"%s\" LIMIT 1", table)
	r, err := db.Query(sql)
	if err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf("\"%s\" does not exist", table)) {
			t.Fatal(err.Error())
		}
		return false
	}
	r.Close()
	return true
}
