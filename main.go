package main

import (
	"fmt"
	"os"

	"github.com/drand/drand/cmd/drand-cli"
)

func main() {
	app := drand.CLI()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("drand: error running app: %s", err)
	}
}
