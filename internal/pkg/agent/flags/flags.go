package flags

import (
	"flag"
	"unicode/utf8"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/env"
)

type Flags interface {
	GetAddr() string
	GetReportInterval() int
	GetPollInterval() int
}

type flags struct {
	addr           string
	reportInterval int
	pollInterval   int
}

func (f *flags) GetAddr() string {
	return f.addr
}

func (f *flags) GetReportInterval() int {
	return f.reportInterval
}

func (f *flags) GetPollInterval() int {
	return f.pollInterval
}

func ParseFlags(env env.Env) Flags {
	f := flags{}

	if env.GetAddr() == nil {
		flag.StringVar(&f.addr, "a", "http://localhost:8080", "http server addr")
	} else {
		f.addr = *env.GetAddr()
	}

	if env.GetReportInterval() == nil {
		flag.IntVar(&f.reportInterval, "r", 10, "report interval")
	} else {
		f.reportInterval = *env.GetReportInterval()
	}

	if env.GetPollInterval() == nil {
		flag.IntVar(&f.pollInterval, "p", 2, "poll interval")
	} else {
		f.pollInterval = *env.GetPollInterval()
	}

	flag.Parse()

	f.addr = normalizeAddr(f.addr)

	return &f
}

func normalizeAddr(addr string) string {
	if utf8.RuneCountInString(addr) < 4 || addr[:4] != "http" {
		return "http://" + addr
	}

	return addr
}
