package main

import (
	"github.com/MaximPolyaev/go-metrics/internal/app/agent"
)

const (
	sPoolInterval   = 2
	sReportInterval = 10
	baseURL         = "http://localhost:8080"
)

func main() {
	agent.Run(sPoolInterval, sReportInterval, baseURL)
}
