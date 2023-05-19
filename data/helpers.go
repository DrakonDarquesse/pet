package data

import (
	"encoding/json"
	"io"
	"net/http"
)

type JsonUtil struct{}

func (j *JsonUtil) ToJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	return e.Encode(data)
}

func (j *JsonUtil) FromJSON(r io.Reader, data any) error {
	e := json.NewDecoder(r)
	return e.Decode(data)
}
