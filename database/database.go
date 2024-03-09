package database

import (
	"database/sql"
	"fmt"
	"os"
)

// Initialize a new database and affect it to DB for all request
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database/dtb.db")
	if err != nil {
		return nil,err
	}
	return db,nil
}
func CreateTables(db *sql.DB) error {
	query, err := os.ReadFile("./database/query.sql")
	if len(query) == 0 {
		return fmt.Errorf("empty file sql")
	}
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	_, err2 := db.Exec(string(query))
	if err2 != nil {
		return fmt.Errorf(err2.Error())
	}
	return nil
}
