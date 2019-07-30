package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type NodeKeyJson struct {
	NodeKey string `json:"node_key"`
}

var NoNodeKeyErr = errors.New("support: no node key found in request")

func ExtractNodeKey(r *http.Request) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// TODO clean up
		log.Println("rip", err)
		return "", NoNodeKeyErr
	}

	input := NodeKeyJson{}
	json.Unmarshal(body, &input)

	if input.NodeKey != "" {
		return input.NodeKey, nil
	} else {
		return "", NoNodeKeyErr
	}
}
