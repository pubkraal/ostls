package data

import (
	"database/sql"
	"log"
)

type Config struct {
	Id      int
	Title   string
	Content string
}

// func ListConfigs(db *sql.DB) []*Config {
//
// }

func LoadConfigByTitle(title string, db *sql.DB) *Config {
	c := &Config{}

	q := db.QueryRow("select id, title, content from config where title=$1", title)
	err := q.Scan(&c.Id, &c.Title, &c.Content)
	if err == sql.ErrNoRows {
		return c
	} else if err != nil {
		log.Println("Error during loading config for ", title, ": ", err)
	}
	return c
}

func (c *Config) Persist(db *sql.DB) {
	var err error
	var rows *sql.Rows
	if c.Id == 0 {
		q := "insert into config (title, content) VALUES ($1, $2);"
		rows, err = db.Query(q, c.Title, c.Content)
	} else {
		q := "update config set title=$1, content=$2 where id=$3;"
		rows, err = db.Query(q, c.Title, c.Content, c.Id)
	}
	defer rows.Close()
	if err != nil {
		log.Println("Error persisting config: ", err)
	}
}
