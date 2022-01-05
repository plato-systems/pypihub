// Package util provides various utilities, including config.
package util

import (
	"io/ioutil"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the global configuration of the server.
var Config struct {
	Server struct {
		Host string
		Port uint16
		TLS  struct {
			Crt string
			Key string
		}
	}

	GitHub struct {
		Owners []string

		Replace []struct {
			Patt uRegexp
			Repl string
		}
	}
}

// LoadConfigFile merges the given TOML configuration file into
// the global configuration.
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
	for _, o := range Config.GitHub.Owners {
		authOwners[o] = true
	}

	return nil
}

const defaultConfig = `
[server]
host = "0.0.0.0"
port = 3141
`

func init() {
	loadConfig([]byte(defaultConfig))
}
