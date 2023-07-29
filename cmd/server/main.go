package main

import "github.com/MaximPolyaev/go-metrics/internal/app/server"

func main() {
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
