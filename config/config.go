// Package config provides the global configuration for the PyPIHub server
package config

import (
	"io/ioutil"

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

func init() {
	toml.Unmarshal(defaultConf, &Conf)
}

func LoadFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return toml.Unmarshal(content, &Conf)
}
