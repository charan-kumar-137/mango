package mango

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data    []byte
	Status  int
	Headers DataMap
	Json    any
}

func (response Response) Send(w http.ResponseWriter) {

	var data []byte

	if response.Json != nil {
		data, _ = json.Marshal(response.Json)
		w.Header().Set("Content-Type", "application/json")
	} else {
		data = response.Data
	}

	if data != nil {
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	}

	for key, values := range response.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(response.Status)
	respWriter := bufio.NewWriter(w)
	respWriter.Write(data)
	respWriter.Flush()
}
