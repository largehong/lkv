package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/largehong/lkv/engine"
	"github.com/largehong/lkv/memkv"
	"github.com/largehong/lkv/processor"
	"github.com/largehong/lkv/watch"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "lkv",
		Usage: "watch remote kv changes and render configuration files in real time",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.yaml",
				Usage:   "configuration file path",
			},
			&cli.StringFlag{
				Name:  "log.level",
				Value: "info",
				Usage: "log level",
			},
			&cli.StringFlag{
				Name:  "templates",
				Value: "./templates",
				Usage: "templates directory",
			},
			&cli.IntFlag{
				Name:    "max",
				Aliases: []string{"m"},
				Value:   100,
				Usage:   "buffer size",
			},
			&cli.IntFlag{
				Name:    "interval",
				Aliases: []string{"i"},
				Value:   3,
				Usage:   "interval",
			},
		},
		Action: func(ctx *cli.Context) error {
			config, err := ParseFromYAML(ctx.String("config"))
			if err != nil {
				return err
			}

			SetConfigFromCLI(ctx, config)

			kv := memkv.New()

			e := engine.New(kv, config.Max, config.Interval)

			w, err := watch.New(config.Watch.Type, config.Watch.Config, config.Watch.Prefixes, e.Callback)
			if err != nil {
				return err
			}

			tpls, err := template.ParseFS(os.DirFS(config.Templates), "*.tpl")
			if err != nil {
				return err
			}
			tpls.Funcs(kv.FuncMaps())

			for _, item := range config.Processors {
				p := processor.New(tpls, item.Src, item.Dst, item.Hook.After)
				for _, prefix := range item.Prefixes {
					e.Register(prefix, p)
				}
			}

			kvs, err := w.Get()
			if err != nil {
				return err
			}
			go e.Run()
			e.Callback(kvs...)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func SetConfigFromCLI(ctx *cli.Context, config *Config) {
	if config.Log.Level == "" {
		config.Log.Level = ctx.String("log.level")
	}

	if config.Interval <= 0 {
		config.Interval = ctx.Int("interval")
	}

	if config.Max <= 0 {
		config.Max = ctx.Int("max")
	}

	if config.Templates == "" {
		config.Templates = ctx.String("templates")
	}
}
