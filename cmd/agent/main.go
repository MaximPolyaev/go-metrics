package main

import (
	"github.com/MaximPolyaev/go-metrics/internal/app/agent"
)

const (
	sPoolInterval   = 2
	sReportInterval = 10
	baseUrl         = "http://localhost:8080"
)

func main() {
	agent.Run(sPoolInterval, sReportInterval, baseUrl)
}
