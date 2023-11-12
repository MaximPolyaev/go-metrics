package handler_test

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/handler"
)

func Example() {
	// create new handler

	h := handler.New(&mockMetricService{})

	// configure routes
	http.Handle("/json/update", h.UpdateByJSONFunc())
	http.Handle("/json/updates", h.BatchUpdateByJSONFunc())
	http.Handle("/json/value", h.GetValueByJSONFunc())

	http.Handle("/simple/ping", h.PingFunc(&sql.DB{}))
	http.Handle("/simple/main", h.MainFunc())
	http.Handle("/simple/update", h.UpdateFunc())
	http.Handle("/simple/value", h.GetValueFunc())

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
