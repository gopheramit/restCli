package main

import (
	"encoding/json"
	"time"

	todo "github.com/gopheramit/todoCLI"
)

type todoResponse struct {
	Result todo.List `json:"results"`
}

func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todo.List `json:"results"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"total_results"`
	}{
		Results:      r.Result,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Result),
	}
	return json.Marshal(resp)

}
