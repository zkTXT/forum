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

func CreateSchema() error {
	stmt := `
		CREATE TABLE IF NOT EXISTS Users (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Username TEXT,
			Email TEXT
		);

		CREATE TABLE IF NOT EXISTS Posts (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT,
			Content TEXT,
			UserID INTEGER,
			FOREIGN KEY (UserID) REFERENCES Users(ID)
		);
	`
	_, err := handler.db.Exec(stmt)
	return err
}

func InsertUser(username, email string) (int64, error) {
	stmt := "INSERT INTO Users (Username, Email) VALUES (?, ?)"
	result, err := handler.db.Exec(stmt, username, email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func QueryUsers() (*sql.Rows, error) {
	stmt := "SELECT ID, Username, Email FROM Users"
	return handler.db.Query(stmt)
}
