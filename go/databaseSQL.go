package forum

import (
	"database/sql"
	"log"
)

type DBHandler struct {
	initialized bool
	db          *sql.DB
}

var handler = new(DBHandler)

func Initialize(dbFile string) {
	if !handler.initialized {
		var err error
		handler.db, err = sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal(err)
		}
		handler.initialized = true
	}
}

func Execute(stmt string) (sql.Result, error) {
	statement, err := handler.db.Prepare(stmt)
	if err != nil {
		log.Fatal(err.Error())
	}
	return statement.Exec()
}

func Query(stmt string) (*sql.Rows, error) {
	return handler.db.Query(stmt)
}
