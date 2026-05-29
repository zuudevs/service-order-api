/**

 filename  : sqlite.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package database

import (
	"os"
	"database/sql"

	_ "modernc.org/sqlite"
)

func initDatabase(db *sql.DB) (*sql.DB, error) {
	file, err := os.ReadFile("./internal/database/schema.sql")

	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(file)); err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectSQLite() (*sql.DB, error) {
	db, err := sql.Open(
		"sqlite",
		"./storage/database.db",
	)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err != nil {
		db, err = initDatabase(db)
		if err != nil {
			_ = db.Close()
			return nil, err
		} 
	}

	return db, nil
}