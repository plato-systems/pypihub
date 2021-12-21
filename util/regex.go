package util

import "regexp"

type uRegexp struct {
	*regexp.Regexp
}

func (r *uRegexp) UnmarshalText(text []byte) error {
	re, err := regexp.Compile(string(text))
	if err != nil {
		return err
	}

	r.Regexp = re
	return nil
}
