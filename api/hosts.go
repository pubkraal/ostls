package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HostsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "{}")
}
