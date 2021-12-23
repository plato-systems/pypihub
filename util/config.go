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

	Owners []string

	Replace []struct {
		Re   uRegexp
		Repl string
	}
}

func LoadConfigFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return loadConfig(content)
}

func loadConfig(text []byte) error {
	if err := toml.Unmarshal(text, &Config); err != nil {
		return err
	}

	authOwners = map[string]bool{}
	for _, o := range Config.Owners {
		authOwners[o] = true
	}

	return nil
}

const defaultConfig = `
host = "0.0.0.0"
port = 3141
`

func init() {
	loadConfig([]byte(defaultConfig))
}
