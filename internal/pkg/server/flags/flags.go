package flags

import (
	"flag"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/env"
)

type Flags interface {
	GetAddr() string
}

type flags struct {
	addr string
}

func (f *flags) GetAddr() string {
	return f.addr
}

func ParseFlags(env env.Env) Flags {
	f := flags{}

	if env.GetAddr() == nil {
		flag.StringVar(&f.addr, "a", ":8080", "http server addr")
		flag.Parse()
	} else {
		f.addr = *env.GetAddr()
	}

	return &f
}
