package config

import (
	env "github.com/SinnerUfa/practicum-metric/internal/env"
	flags "github.com/SinnerUfa/practicum-metric/internal/flags"
)

func Load(v any, args []string) error {
	if err := flags.Load(v, args); err != nil {
		return err
	}
	if err := env.Load(v); err != nil {
		return err
	}
	return nil
}
