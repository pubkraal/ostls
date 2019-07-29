package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func stub(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "POST" {
		log.Println("Request POST request.")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body")
		} else {
			log.Println(string(body))
		}
	}

	fmt.Fprintf(w, "This is a stub call, it doesn't work.")
}

func main() {
	cert := flag.String("cert", "./cert.pem", "valid TLS certificate for hosting https")
	key := flag.String("key", "./key.pem", "valid TLS Private key for hosting https")
	port := flag.Int("port", 0, "the port to host on. Leave empty to default to 80 or 443 depending on TLS config")
	host := flag.String("host", "", "the host to listen on, leave empty for all interfaces.")
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

	router := httprouter.New()

	// Normal route
	router.GET("/", stub)

	// User routes, stuff like management should be in these.

	// API routes for the client taken from the facebook test http
	// server. Implementation will require more attention obviously.
	// https://github.com/osquery/osquery/blob/master/tools/tests/test_http_server.py
	router.GET("/api/config", stub)
	router.POST("/api/enroll", stub)
	router.POST("/api/log", stub)
	router.POST("/api/distributed-read", stub)
	router.POST("/api/distributed-write", stub)
	router.POST("/api/test-read-requests", stub)
	router.POST("/api/carve-init", stub)
	router.POST("/api/carve-block", stub)

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
		log.Fatal(http.ListenAndServeTLS(hostname, *cert, *key, nil))
	} else {
		log.Fatal(http.ListenAndServe(hostname, router))
	}
}
