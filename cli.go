package main

import (
	"errors"
	"github.com/ElaraLang/elara/base"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "Elara",
		Usage: "ExecuteFull Elara Code",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "script",
				Value: false,
				Usage: "Script Mode (print the result of every expression)",
			},
		},
		Action: func(c *cli.Context) error {
			fileName := c.Args().Get(0)
			if fileName == "" {
				return errors.New("no file provided to execute - nothing to do")
			}

			scriptMode := c.Bool("script")
			base.ExecuteFull(fileName, scriptMode)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
