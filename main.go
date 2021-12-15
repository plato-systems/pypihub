package main

import (
	"fmt"

	"github.com/plato-systems/pypihub/config"
)

func main() {
	err := config.Load("pypihub.toml")
	fmt.Println(err, config.Conf)
}
