package config

import "regexp"

type regex struct {
	*regexp.Regexp
}

func (r *regex) UnmarshalText(text []byte) error {
	re, err := regexp.Compile(string(text))
	if err != nil {
		return err
	}

	r.Regexp = re
	return nil
}
