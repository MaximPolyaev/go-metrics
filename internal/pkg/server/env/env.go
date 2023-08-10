package env

import (
	cenv "github.com/caarlos0/env/v9"
)

type Env interface {
	GetAddr() *string
}

type env struct {
	Addr *string `env:"ADDRESS"`
}

func (e *env) GetAddr() *string {
	return e.Addr
}

func ParseEnv() (Env, error) {
	var e env

	err := cenv.Parse(&e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
