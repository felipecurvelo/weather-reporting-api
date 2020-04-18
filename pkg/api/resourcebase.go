package api

import (
	"encoding/json"
	"net/http"
)

type ResourceBase struct {
}

func (r *ResourceBase) SetResponse(status int, response interface{}, w http.ResponseWriter) {
	b := response

	_, ok := response.([]byte)
	if !ok {
		b, _ = json.Marshal(response)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b.([]byte))
}
