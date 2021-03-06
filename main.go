package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/pubkraal/ostls/api"
	"github.com/pubkraal/ostls/util"
)

type Logger struct {
	handler http.Handler
}

func (l Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Proto, r.Method, r.URL, r.Header.Get("User-Agent"))
	l.handler.ServeHTTP(w, r)
}

func stub(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "{}")
}

func stubWriter(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body")
		} else {
			log.Println(string(body))
		}
	}

	fmt.Fprintf(w, "{}")
}

func main() {
	defaultDb := util.FirstNonEmpty(os.Getenv("OSTLS_DSN"), "postgres://localhost/ostls")
	defaultEs := util.FirstNonEmpty(os.Getenv("OSTLS_ES"), "https://localhost:9200/")
	defaultCert := util.FirstNonEmpty(os.Getenv("OSTLS_CERT"), "./cert.pem")
	defaultKey := util.FirstNonEmpty(os.Getenv("OSTLS_KEY"), "./key.pem")
	defaultPort, _ := strconv.Atoi(os.Getenv("OSTLS_PORT"))
	defaultHost := os.Getenv("OSTLS_HOST")

	// Flags for external services
	dsn := flag.String("dsn", defaultDb, "the database DSN")
	es := flag.String("es", defaultEs, "the connectstring for elasticsearch")

	// Flags for tls/non-tls
	cert := flag.String("cert", defaultCert, "valid TLS certificate for hosting https")
	key := flag.String("key", defaultKey, "valid TLS Private key for hosting https")
	port := flag.Int("port", defaultPort, "the port to host on. Leave empty to default to 80 or 443 depending on TLS config")
	host := flag.String("host", defaultHost, "the host to listen on, leave empty for all interfaces.")

	// Flags for enrollment management
	// secret := flag.String("secret", "", "the file in which the shared secret is stored.")
	flag.Parse()

	if *port < 0 || *port > 65535 {
		log.Fatalf("Invalid port number: %d\n", *port)
	}

	tlsReady := true
	if _, err := os.Stat(*cert); err != nil {
		log.Println("Invalid certificate", err)
		tlsReady = false
	}
	if _, err := os.Stat(*key); err != nil {
		log.Println("Invalid private key", err)
		tlsReady = false
	}

	dbHandle := api.InitializeDatabase(*dsn)
	defer dbHandle.Close()
	api.InitializeElastic(*es)

	router := httprouter.New()

	// Normal route
	router.GET("/", stub)

	// User routes, stuff like management should be in these.

	// API routes for the client taken from the facebook test http
	// server. Implementation will require more attention obviously.
	// https://github.com/osquery/osquery/blob/master/tools/tests/test_http_server.py
	router.POST("/api/config", api.Config)
	router.POST("/api/enroll", api.Enroll)
	router.POST("/api/log", api.AcceptLog)
	router.POST("/api/distributed-read", stub)
	router.POST("/api/distributed-write", stub)

	// Not implemented
	router.GET("/api/config", stub)
	router.POST("/api/test-read-requests", stub)
	router.POST("/api/carve-init", stub)
	router.POST("/api/carve-block", stub)

	router.GET("/api/hosts", api.HostsList)

	if *port == 0 {
		if tlsReady {
			*port = 443
		} else {
			*port = 80
		}
	}

	hostname := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Listening on", hostname)

	if tlsReady {
		log.Fatal(http.ListenAndServeTLS(hostname, *cert, *key, Logger{router}))
	} else {
		log.Fatal(http.ListenAndServe(hostname, Logger{router}))
	}
}
