package data

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Host struct {
	Id         int
	Identifier string
	UUID       uuid.UUID
	Hostname   string
	Token      uuid.UUID
	Enrolled   time.Time
}

func LoadHost(id int, db *sql.DB) *Host {
	h := &Host{}

	q := db.QueryRow("select id, identifier, uuid, hostname, token, enrolled from host where id = $1", id)
	err := q.Scan(&h.Id, &h.Identifier,
		&h.UUID, &h.Hostname,
		&h.Token, &h.Enrolled)
	if err == sql.ErrNoRows {
		return h
	} else if err != nil {
		log.Println("Error during loading host. ", err)
		return h
	}

	return h
}

func LoadHostByUUID(uuid string, db *sql.DB) *Host {
	h := &Host{}

	q := db.QueryRow("select id, identifier, uuid, hostname, token, enrolled from host where uuid = $1", uuid)
	err := q.Scan(&h.Id, &h.Identifier,
		&h.UUID, &h.Hostname,
		&h.Token, &h.Enrolled)
	if err == sql.ErrNoRows {
		return h
	} else if err != nil {
		log.Println("Error during loading host. ", err)
		return h
	}

	return h
}

func ListValidTokens(db *sql.DB) []uuid.UUID {
	var uuids []uuid.UUID
	rows, err := db.Query("select token from host")
	if err != nil {
		return uuids
	}
	defer rows.Close()

	for rows.Next() {
		var localUuid uuid.UUID
		err := rows.Scan(&localUuid)
		if err != nil {
			continue
		}

		uuids = append(uuids, localUuid)
	}

	return uuids
}

func (h *Host) Persist(db *sql.DB) {
	if h.Id == 0 {
		q := "insert into host (identifier, uuid, hostname, token, enrolled) VALUES ($1, $2, $3, $4, $5)"

		rows, err := db.Query(q, h.Identifier, h.UUID, h.Hostname, h.Token, h.Enrolled)
		defer rows.Close()
		if err != nil {
			log.Println("Error storing new host: ", err)
		}
	} else {
		// Update
		q := "update host set identifier=$1, hostname=$2, token=$3 where uuid=$4"

		rows, err := db.Query(q, h.Identifier, h.Hostname, h.Token, h.UUID)
		defer rows.Close()
		if err != nil {
			log.Println("Error updating host: ", err)
		}
	}
}
