package server

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/handler"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/mem_storage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/storage"
	"net/http"
)

const (
	updateAction = "/update/"

	updateGaugeAction   = updateAction + metric.GaugeType + "/"
	updateCounterAction = updateAction + metric.CounterType + "/"
)

type middleware func(http.Handler) http.Handler

func Run(addr string) error {
	s := mem_storage.NewMemStorage()

	mux := createServeMux(s)

	return http.ListenAndServe(addr, mux)
}

func createServeMux(s storage.Storage) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())

	mux.Handle(
		updateGaugeAction,
		conveyor(
			http.StripPrefix(updateGaugeAction, handler.GaugeFunc(s)),
			allowedMethodMiddleware,
		),
	)

	mux.Handle(
		updateCounterAction,
		conveyor(
			http.StripPrefix(updateCounterAction, handler.CounterFunc(s)),
			allowedMethodMiddleware,
		),
	)

	mux.Handle(
		updateAction,
		conveyor(http.HandlerFunc(handler.IncorrectMetric), allowedMethodMiddleware),
	)

	return mux
}

func conveyor(h http.Handler, middlewares ...middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}

	return h
}

func allowedMethodMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method "+r.Method+" not allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
