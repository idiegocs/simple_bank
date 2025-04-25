package db

import (
	"database/sql"
	"log"
	"testing"
	"os"
	 _"github.com/lib/pq"

)

var testQueries *Queries
var testDB *sql.DB 

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:dics@localhost:5532/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {

	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	log.Println("Connected ",testQueries)

	os.Exit(m.Run())

	

}
