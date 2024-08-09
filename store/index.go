package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var store *sql.DB

func InitDatabase(path string) error {
	var err error
	store, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	// databases init

	_, err = store.Exec(`

    CREATE TABLE IF NOT EXISTS clipboard_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        type TEXT,
        categories TEXT

    );`)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return err
	}
	return nil
}
