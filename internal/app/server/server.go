package server

import "net/http"

type middleware func(http.Handler) http.Handler

func Run(addr string) error {
	mux := createServeMux()

	return http.ListenAndServe(addr, mux)
}

func createServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())

	mux.Handle(
		updateGaugeAction,
		conveyor(
			http.StripPrefix(updateGaugeAction, http.HandlerFunc(gaugeHandler)),
			allowedMethodMiddleware,
		),
	)

	mux.Handle(
		updateCounterAction,
		conveyor(
			http.StripPrefix(updateCounterAction, http.HandlerFunc(counterHandler)),
			allowedMethodMiddleware,
		),
	)

	mux.Handle(
		updateAction,
		conveyor(http.HandlerFunc(incorrectMetricHandler), allowedMethodMiddleware),
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
