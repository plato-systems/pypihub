// Package util provides various utilities, including config
package util

import (
	"io/ioutil"

	"github.com/pelletier/go-toml/v2"
)

var Config struct {
	Host string
	Port uint16
	TLS  struct {
		Cert string
		Key  string
	}

	Replace []struct {
		Re   uRegexp
		Repl string
	}
}

func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return toml.Unmarshal(content, &Config)
}

var defaultConfig = []byte(`
host = "0.0.0.0"
port = 3141
`)

func init() {
	toml.Unmarshal(defaultConfig, &Config)
}
