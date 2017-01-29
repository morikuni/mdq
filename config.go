package mdq

import (
	"database/sql"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DBs []DBConfig `yaml:"dbs"`
}

type DBConfig struct {
	Name   string   `yaml:"name"`
	Driver string   `yaml:"driver"`
	DSN    string   `yaml:"dsn"`
	Tags   []string `yaml:"tags,flow"`
}

func CreateDBsFromConfig(r io.Reader, tag string) ([]DB, error) {
	conf, err := ParseConfig(r)
	if err != nil {
		return nil, err
	}

	var dbs []DB
	for _, dbc := range conf.DBs {
		if tag != "" {
			if len(dbc.Tags) == 0 {
				continue
			}
			for _, t := range dbc.Tags {
				if tag != t {
					continue
				}
			}
		}
		con, err := sql.Open(dbc.Driver, dbc.DSN)
		if err != nil {
			panic(err)
		}
		dbs = append(dbs, NewDB(dbc.Name, con))
	}

	return dbs, nil
}

func ParseConfig(r io.Reader) (Config, error) {
	bs, err := ioutil.ReadAll(r)
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
