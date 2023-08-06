package flags

import "flag"

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

func ParseFlags() Flags {
	f := flags{}

	flag.StringVar(&f.addr, "a", "localhost:8080", "http server addr")
	flag.IntVar(&f.reportInterval, "r", 10, "report interval")
	flag.IntVar(&f.pollInterval, "p", 2, "poll interval")

	flag.Parse()

	return &f
}
