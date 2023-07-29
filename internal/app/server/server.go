package server

import "net/http"

type middleware func(http.Handler) http.Handler

func Run(addr string) error {
	mux := createServeMux()

	return http.ListenAndServe(addr, mux)
}

func createServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", conveyor(createEmptyHandlerFunc(), allowedMethodMiddleware))

	mux.Handle(
		updateGaugeAction,
		conveyor(http.HandlerFunc(gaugeHandler), allowedMethodMiddleware),
	)

	mux.Handle(
		updateCounterAction,
		conveyor(http.HandlerFunc(counterHandler), allowedMethodMiddleware),
	)

	return mux
}

func conveyor(h http.Handler, middlewares ...middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}

	return h
}

func createEmptyHandlerFunc() http.HandlerFunc {
	return func(_ http.ResponseWriter, _ *http.Request) {}
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
