package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var dbHandle *sql.DB

func InitializeDatabase(dsn string) *sql.DB {
	localDBHandle, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening the database", err)
	}
	err = localDBHandle.Ping()
	if err != nil {
		log.Fatal("Error pinging the database", err)
	}

	dbHandle = localDBHandle

	return localDBHandle
}

func writeFailure(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"node_invalid\": true}")
}
