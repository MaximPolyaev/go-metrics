package server

import (
	"io"
	"net/http"
)

func gaugeHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, "gauge"); err != nil {
		panic(err)
	}
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, "counter"); err != nil {
		panic(err)
	}
}
