package storage

import (
	dsql "database/sql"
	"net/url"
	//Since this is the place where the database type selection occous
	//it seems to be the best place for this blank import
	_ "github.com/go-sql-driver/mysql"

	"github.com/EmpregoLigado/code-challenge/crypt"
	"github.com/EmpregoLigado/code-challenge/storage/interfaces"
	"github.com/EmpregoLigado/code-challenge/storage/memory"
	"github.com/EmpregoLigado/code-challenge/storage/sql"
)

//NewJob instantiates a job database based on the db url
//Currently only memory and mysql/mariadb are supported
func NewJob(cipher crypt.Cipher, dburl string) (interfaces.Job, error) {
	u, err := url.Parse(dburl)
	if err != nil {
		return nil, err
	}
	dbtype := u.Scheme
	u.Scheme = ""

	switch dbtype {
	case "mysql":
		q := u.Query()
		q.Set("parseTime", "true")
		u.RawQuery = q.Encode()
		db, err := dsql.Open(dbtype, u.String()[2:])
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err != nil {
			return nil, err
		}
		return sql.NewJob(cipher, db), nil
	default:
		return memory.NewJob(cipher), nil
	}
}
