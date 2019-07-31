package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/elastic/go-elasticsearch/v7"
)

var dbHandle *sql.DB
var esHandle *elasticsearch.Client

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

func InitializeElastic(es string) {
	esUrl, err := url.Parse(es)
	if err != nil {
		log.Fatal("Elasticsearch string misconfigured.", err)
	}

	esAddress := fmt.Sprintf("%s://%s", esUrl.Scheme, esUrl.Host)
	esUser := esUrl.User.Username()
	esPass, _ := esUrl.User.Password()

	cfg := elasticsearch.Config{
		Addresses: []string{
			esAddress,
		},
		Username: esUser,
		Password: esPass,
	}
	localEs, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to create ElasticSearch client. ", err)
	}

	_, err = localEs.Info()
	if err != nil {
		log.Fatal("Could not connect to elasticsearch. ", err)
	}
	esHandle = localEs
}

func writeFailure(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"node_invalid\": true}")
}
