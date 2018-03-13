package sql

import (
	"database/sql"
	"net/url"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/EmpregoLigado/code-challenge/crypt/null"
	"github.com/EmpregoLigado/code-challenge/storage/test"
)

func TestJobBackend(t *testing.T) {
	backends := []test.Backend{}
	cipher := &null.Cipher{}
	dburl := os.Getenv("DB")
	if dburl == "" {
		t.Skip("skipping DB url not found on environment")
	}

	u, err := url.Parse(dburl)
	if err != nil {
		t.Fatal("unable to parse db url")
	}
	dbtype := u.Scheme

	u.Scheme = ""
	q := u.Query()
	q.Set("parseTime", "true")
	u.RawQuery = q.Encode()

	db, err := sql.Open(dbtype, u.String()[2:])
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	backends = append(backends, test.Backend{Name: "mysql", Backend: NewJob(cipher, db)})
	test.JobBackendTest(t, backends)
}
