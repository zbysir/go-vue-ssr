package main

import (
	"github.com/bysir-zl/vue-ssr/pkg/vuessr"
	"github.com/urfave/cli"
	"os"
)

func main() {
	c := cli.NewApp()
	c.Name = "vuessr"
	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "src",
			Usage: "the .vue file dir",
		},
		&cli.StringFlag{
			Name:  "to",
			Value: "./internal/vuetpl",
			Usage: "genera code dir (default is ./internal/vuetpl)",
		},
	}

	c.Action = func(c *cli.Context) (err error) {
		src := c.String("src")
		if src == "" {
			panic("invalid src")
		}
		to := c.String("to")
		err = vuessr.GenAllFile(src, to)
		if err != nil {
			return
		}

		return
	}

	err := c.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
