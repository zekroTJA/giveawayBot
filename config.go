package main

import (
	"io/ioutil"
	"github.com/go-yaml/yaml"
)


type ConfigData struct {
	Token      string   `yaml:"token"`
	Emote      string   `yaml:"emote"`
	Authorized []string `yaml:"authorized"`
}

type Config struct {
	Path string
	Data *ConfigData
}

func NewConfig(path string) (*Config, error) {
	b_data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := ConfigData{}
	err = yaml.Unmarshal(b_data, &data)
	if err != nil {
		return nil, err
	}
	c := &Config{ path, &data }
	return c, nil
}