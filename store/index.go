package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var store *sql.DB

func ShutDownDatabase() {
	if store != nil {
		err := store.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Database connection closed")
}

func InitDatabase(path string) error {
	var err error
	store, err = sql.Open("sqlite3", path) // Changed "sqlite" to "sqlite3" to match the import
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
		return err
	}

	// Close the database connection when the program exits
	// defer func() {
	// 	if err := store.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	log.Println("Database is connected successfully")
	return nil
}
