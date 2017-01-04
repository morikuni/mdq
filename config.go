package mdq

import (
	"database/sql"
	"io/ioutil"
	"regexp"

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

func CreateDBsFromFile(path string, filter *regexp.Regexp) ([]DB, error) {
	conf, err := ParseFile(path)
	if err != nil {
		return nil, err
	}

	var dbs []DB
	for _, dbc := range conf.DBs {
		if filter != nil && !filter.MatchString(dbc.Name) {
			continue
		}
		con, err := sql.Open(dbc.Driver, dbc.DSN)
		if err != nil {
			panic(err)
		}
		dbs = append(dbs, NewDB(dbc.Name, con))
	}

	return dbs, nil
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
