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
func MatchGQLParam(field, param string, q GraphQLRequest) (string, error) {
	m := regexp.MustCompile(fmt.Sprintf(
		`%s\(.*%s:\s*(?:\$(\w+)|"([^"]*)").*\){`,
		field, param,
	)).FindStringSubmatch(q.Query)
	if m == nil {
		return "", fmt.Errorf("no match")
	}
	if m[1] == "" {
		return m[2], nil
	}

	res, ok := q.Variables[m[1]].(string)
	if !ok {
		return "", fmt.Errorf("non-string variable")
	}
	return res, nil
}
