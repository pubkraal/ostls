package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func stub(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "This is a stub call, it doesn't work.")
}

func main() {
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

	log.Fatal(http.ListenAndServe(":8080", router))
}
