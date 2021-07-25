package connection

import "database/sql"

var db *sql.DB

func GetConnection() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	db, err = sql.Open("sqlite3", "./db/tibiawiki.db")
	if err != nil {
		println("error coneccion")
		panic(err)

	}
	return db
}
