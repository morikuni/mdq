package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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

	home := os.Getenv("HOME")

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
	dbs := make(map[string]mdq.DB)
	for _, dbc := range conf.DBs {
		con, err := sql.Open(dbc.Driver, dbc.DSN)
		if err != nil {
			panic(err)
		}
		dbs[dbc.Name] = mdq.NewDB(con)
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
