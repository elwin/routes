package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			stravaPull,
			poster,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
