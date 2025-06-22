package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/largehong/lkv/engine"
	"github.com/largehong/lkv/memkv"
	"github.com/largehong/lkv/processor"
	"github.com/largehong/lkv/watch"
	_ "github.com/largehong/lkv/watch/etcdv3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func GetLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func GetLogFormat(f string) logrus.Formatter {
	switch strings.ToLower(f) {
	case "json":
		return &logrus.JSONFormatter{}
	case "text":
		return &logrus.TextFormatter{DisableColors: true}
	default:
		return &logrus.TextFormatter{}
	}
}

func main() {
	app := &cli.App{
		Name:  "lkv",
		Usage: "watch remote kv changes and render configuration files in real time",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.yml",
				Usage:   "configuration file path",
			},
			&cli.StringFlag{
				Name:  "log.level",
				Value: "info",
				Usage: "log level",
			},
			&cli.StringFlag{
				Name:  "log.format",
				Value: "text",
				Usage: "log format, json or text",
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
			&cli.BoolFlag{
				Name:  "once",
				Usage: "run once",
				Value: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			config, err := ParseFromYAML(ctx.String("config"))
			if err != nil {
				return err
			}

			SetConfigFromCLI(ctx, config)

			logrus.SetLevel(GetLogLevel(config.Log.Level))
			logrus.SetFormatter(GetLogFormat(config.Log.Format))

			kv := memkv.New()

			e := engine.New(kv, config.Max, config.Interval)

			tpls, err := template.New("lkv").Funcs(kv.FuncMaps()).Funcs(processor.FuncMaps()).ParseFS(os.DirFS(config.Templates), "*.tpl")
			if err != nil {
				return err
			}

			w, err := watch.New(config.Watch.Type, config.Watch.Config, config.Watch.Prefixes, e.Callback)
			if err != nil {
				return err
			}

			for _, item := range config.Processors {
				tpl := tpls.Lookup(item.Src)
				if tpl == nil {
					return errors.New("not found template: " + item.Src)
				}
				p := processor.New(tpl, item.Src, item.Dst, item.Hook.After)
				for _, prefix := range item.Prefixes {
					e.Register(prefix, p)
				}
			}

			kvs, err := w.Get()
			if err != nil {
				return err
			}

			if ctx.Bool("once") {
				e.Once(kvs...)
				return nil
			}
			go e.Run()
			e.Callback(kvs...)
			select {}
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

	if config.Log.Format == "" {
		config.Log.Format = ctx.String("log.format")
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
