package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"text/template"

	"gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
	"github.com/morikuni/mdq"
)

func main() {
	var format string
	var query string
	var config string
	var silent bool
	var target string

	home := os.Getenv("HOME")

	flag.StringVar(&target, "target", "", "target")
	flag.StringVar(&format, "f", "", "format")
	flag.StringVar(&query, "q", "", "query")
	flag.StringVar(&config, "c", home+"/.config/mdq/config.yaml", "config")
	flag.BoolVar(&silent, "s", false, "print error")
	flag.Parse()

	if query == "" {
		panic("query is empty")
	}

	if config == "" {
		panic("config is empty")
	}

	targetReg, err := regexp.Compile(target)
	if err != nil {
		panic(err)
	}

	reporter := mdq.DefaultReporter
	if silent {
		reporter = mdq.SilentReporter
	}

	bs, err := ioutil.ReadFile(config)
	if err != nil {
		panic(err)
	}
	var conf mdq.Config
	err = yaml.Unmarshal(bs, &conf)
	if err != nil {
		panic(err)
	}
	var dbs []mdq.DB
	for _, dbc := range conf.DBs {
		if !targetReg.MatchString(dbc.Name) {
			continue
		}
		con, err := sql.Open(dbc.Driver, dbc.DSN)
		if err != nil {
			panic(err)
		}
		dbs = append(dbs, mdq.NewDB(dbc.Name, con))
	}

	cluster := mdq.NewCluster(dbs, reporter)

	results := cluster.Query(query)

	if format != "" {
		t := template.New("sql")
		t, err = t.Parse(format)
		if err != nil {
			panic(err)
		}
		t.Execute(os.Stdout, results)
		return
	}

	json, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
}

type Result struct {
	Headers []string
	Rows    []map[string]interface{}
}
