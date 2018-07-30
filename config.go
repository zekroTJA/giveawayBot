package main

import (
	"io/ioutil"
	"github.com/go-yaml/yaml"
)

// ConfigFile contains the parsed data values
// of the config file
type ConfigData struct {
	Token      string   `yaml:"token"`
	Prefix     string   `yaml:"prefix"`
	Admin 	   string   `yaml:"admin"`
	Emote      string   `yaml:"emote"`
	Language   string   `yaml:"language"`
}

// Config contains the ConfigData and the path of
// teh config file
type Config struct {
	Path string
	Data *ConfigData
}

// NewConfig creates a new instance of Config
// getting file name + path of the config file passed
// as argument and returning the Config instance and error.
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