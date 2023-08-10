package env

import (
	cenv "github.com/caarlos0/env/v9"
)

type Env interface {
	GetAddr() *string
	GetReportInterval() *int
	GetPollInterval() *int
}

type env struct {
	Addr           *string `env:"ADDRESS"`
	ReportInterval *int    `env:"REPORT_INTERVAL"`
	PollInterval   *int    `env:"POLL_INTERVAL"`
}

func (e *env) GetAddr() *string {
	return e.Addr
}

func (e *env) GetReportInterval() *int {
	return e.ReportInterval
}

func (e *env) GetPollInterval() *int {
	return e.PollInterval
}

func ParseEnv() (Env, error) {
	var e env

	err := cenv.Parse(&e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
