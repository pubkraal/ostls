package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/julienschmidt/httprouter"
)

type LogRequest struct {
	NodeKey string        `json:"node_key"`
	LogType string        `json:"log_type"`
	Data    []interface{} `json:"data"`
}

func AcceptLog(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body")
			fmt.Fprintf(w, "{}")
			return
		}

		// Deserialize request
		logReq := LogRequest{}
		err = json.Unmarshal(body, &logReq)
		if err != nil {
			log.Println("Error unmarshalling: ", err)
			fmt.Fprintf(w, "{}")
			return
		}

		// TODO: verify node key

		// Create one message for ES per blob item
		var dataBlob []byte
		for idx := range logReq.Data {
			dataBlob, err = json.Marshal(logReq.Data[idx])
			if err != nil {
				log.Println("Couldn't marshal json blob. ", err)
				continue
			}

			req := esapi.IndexRequest{
				Index: logReq.LogType,
				Body:  strings.NewReader(string(dataBlob)),
			}

			_, err = req.Do(context.Background(), esHandle)
			if err != nil {
				log.Println("Error indexing log. ", err)
			}
		}

	}

	fmt.Fprintf(w, "{}")
}
