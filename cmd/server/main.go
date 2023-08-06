package main

import "github.com/MaximPolyaev/go-metrics/internal/app/server"

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
