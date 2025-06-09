package response_types

import (
	"encoding/json"
	"net/http"
)

type ResponseBody[T any] struct {
	Data  T       `json:"data"`
	Error *string `json:"error"`
}

func Write[T any](w http.ResponseWriter, httpCode int, err error, data T) {
	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json")

	var r ResponseBody[T]

	r.Data = data
	if err != nil {
		msg := err.Error()
		r.Error = &msg
	}

	b, _ := json.Marshal(r)
	w.Write(b)
}

func WriteErrorNoBody(w http.ResponseWriter, httpCode int, err error) {
	Write[any](w, httpCode, err, nil)
}

func WriteOkEmptyJsonBody(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func WriteOkJsonBody[T any](w http.ResponseWriter, body T) {
	WriteJsonBody(w, http.StatusOK, body)
}

func WriteJsonBody[T any](w http.ResponseWriter, httpStatusCode int, body T) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(httpStatusCode)
	var r ResponseBody[T]
	r.Data = body
	r.Error = nil
	b, _ := json.Marshal(r)
	w.Write(b)
}
