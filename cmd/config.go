package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Log struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Watch struct {
	Type     string   `yaml:"type"`
	Config   any      `yaml:"config"`
	Prefixes []string `yaml:"prefixes"`
}

type Hook struct {
	After string `yaml:"after"`
}

type Processor struct {
	Src      string   `yaml:"src"`
	Dst      string   `yaml:"dst"`
	Prefixes []string `yaml:"prefixes"`
	Hook     Hook     `yaml:"hook"`
}

type Config struct {
	Log        Log         `yaml:"log"`
	Interval   int         `yaml:"interval"`
	Max        int         `yaml:"max"`
	Templates  string      `yaml:"templates"`
	Watch      Watch       `yaml:"watch"`
	Processors []Processor `yaml:"processors"`
}

func ParseFromYAML(filename string) (config *Config, err error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &config)
	return
}
