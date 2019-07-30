package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type EnrollResponse struct {
	NodeKey uuid.UUID `json:"node_key"`
}

func Enroll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("rip", err)
	}

	log.Println("Got a new machine in!", string(body))

	nodeKey := uuid.Must(uuid.NewRandom())
	response := &EnrollResponse{NodeKey: nodeKey}
	fmt.Println(response)

	// Now marshall and write the json out again
	resp, err := json.Marshal(response)
	if err != nil {
		log.Println("Rip!", err)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}
