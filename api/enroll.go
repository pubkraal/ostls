package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pubkraal/ostls/data"
)

type HostDetails struct {
	SystemInfo SystemInfo `json:"system_info"`
}
type SystemInfo struct {
	LocalHostname string `json:"local_hostname"`
	UUID          string `json:"uuid"`
}

type EnrollRequest struct {
	EnrollSecret   string      `json:"enroll_secret"`
	HostIdentifier string      `json:"host_identifier"`
	PlatformType   string      `json:"platform_type"`
	HostDetails    HostDetails `json:"host_details"`
}

type EnrollResponse struct {
	NodeKey uuid.UUID `json:"node_key"`
}

func Enroll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("rip", err)
	}

	input := EnrollRequest{}
	json.Unmarshal(body, &input)

	nodeKey := uuid.Must(uuid.NewRandom())

	// Load host or generate new one
	host := data.LoadHostByUUID(input.HostDetails.SystemInfo.UUID, dbHandle)
	// Update machine with new information!
	host.Token = nodeKey
	host.Identifier = input.HostIdentifier
	host.UUID, err = uuid.Parse(input.HostDetails.SystemInfo.UUID)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid UUID")
		return
	}
	host.Hostname = input.HostDetails.SystemInfo.LocalHostname
	if host.Id == 0 {
		host.Enrolled = time.Now()
	}

	host.Persist(dbHandle)

	response := &EnrollResponse{NodeKey: nodeKey}

	// Now marshall and write the json out again
	resp, err := json.Marshal(response)
	if err != nil {
		log.Println("Rip!", err)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}
