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
	"database/sql"

	_ "modernc.org/sqlite"
)

func ConnectSQLite() (*sql.DB, error) {
	db, err := sql.Open(
		"sqlite",
		"./storage/database.db",
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}