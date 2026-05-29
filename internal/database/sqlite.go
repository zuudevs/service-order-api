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
	"path/filepath"

	_ "modernc.org/sqlite"
)

func initDatabase(
	db *sql.DB, 
	wd string,
) (*sql.DB, error) {
	file, err := os.ReadFile(filepath.Join(wd, "internal", "database", "schema.sql"))

	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(file)); err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectSQLite() (*sql.DB, error) {
	wd, err := os.Getwd()

	if err != nil {
        return nil, err
    }

	storageDir := filepath.Join(wd, "storage")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

    db, err := sql.Open("sqlite", filepath.Join(wd, "storage", "database.db"))
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return initDatabase(db, wd)
}