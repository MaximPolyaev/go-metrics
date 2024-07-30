// Package config for configure services
// Package will be parsing configs
package config

import (
	"time"
)

func convStrIntervalToInt(interval string) (int, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return 0, err
	}

	return int(d / time.Second), nil
}
