// Package config provides the global configuration for the PyPIHub server
package config

import (
	"io/ioutil"
	"regexp"

	"github.com/pelletier/go-toml/v2"
)

var Conf struct {
	Host string
	Port uint16
	TLS  struct {
		Cert string
		Key  string
	}

	Replace []struct {
		Re   regex
		Repl string
	}
}

var defaultConf = []byte(`
host = "0.0.0.0"
port = 3141
`)

func Load(filename string) error {
	toml.Unmarshal(defaultConf, &Conf)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return toml.Unmarshal(content, &Conf)
}

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
