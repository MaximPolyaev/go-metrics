package flags

import "flag"

type Flags interface {
	GetAddr() string
}

type flags struct {
	addr string
}

func (f *flags) GetAddr() string {
	return f.addr
}

func ParseFlags() Flags {
	f := flags{}

	flag.StringVar(&f.addr, "a", ":8080", "http server addr")

	flag.Parse()

	return &f
}
