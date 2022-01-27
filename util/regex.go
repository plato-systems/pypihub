package util

import (
	"fmt"
	"regexp"
)

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

// MatchGQLParam extracts a param included in a GraphQL query string.
func MatchGQLParam(field, param, query string) []string {
	return regexp.MustCompile(fmt.Sprintf(
		`%s\(.*%s:\s*(?:\$(\w+)|"([^"]*)").*\){`,
		field, param,
	)).FindStringSubmatch(query)
}
