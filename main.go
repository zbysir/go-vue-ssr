package main

import (
	"github.com/urfave/cli/v2"
	"github.com/zbysir/go-vue-ssr/internal/pkg/log"
	"github.com/zbysir/go-vue-ssr/internal/version"
	"github.com/zbysir/go-vue-ssr/pkg/vuessr"
	"os"
)

func main() {
	c := cli.NewApp()
	c.Name = "go-vue-ssr"
	c.Description = "Hey vue go"
	c.Version = version.Version
	c.Usage = "Vue to Go compiler"

	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "src",
			Usage: "The .vue files dir",
		},
		&cli.StringFlag{
			Name:  "to",
			Value: "./internal/vuetpl",
			Usage: "Dist dir",
		},
		&cli.StringFlag{
			Name:  "pkg",
			Usage: "pkg name",
		},
	}

	c.Action = func(c *cli.Context) (err error) {
		src := c.String("src")
		if src == "" {
			panic("invalid src")
		}
		to := c.String("to")
		pkg := c.String("pkg")
		err = vuessr.GenAllFile(src, to, pkg)
		if err != nil {
			return
		}

		return
	}

	err := c.Run(os.Args)
	if err != nil {
		log.Errorf("%v", err)
	}
}
