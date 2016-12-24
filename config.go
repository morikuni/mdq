package mdq

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DBs []DBConfig `yaml:"dbs"`
}

type DBConfig struct {
	Name   string `yaml:"name"`
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

func ParseFile(path string) (Config, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var conf Config
	err = yaml.Unmarshal(bs, &conf)
	if err != nil {
		return Config{}, err
	}
	return conf, nil
}
