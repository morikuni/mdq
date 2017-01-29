package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/morikuni/mdq"
	"github.com/spf13/pflag"
)

func main() {
	home := os.Getenv("HOME")

	flag := pflag.NewFlagSet("mdq", pflag.ContinueOnError)
	tag := flag.String("tag", "", "database tag")
	format := flag.String("format", "", "golang template string")
	query := flag.StringP("query", "q", "", "SQL")
	config := flag.String("config", home+"/.config/mdq/config.yaml", "path to config file")
	silent := flag.Bool("silent", false, "ignore errors from databases")

	flag.Parse(os.Args[1:])

	if *query == "" {
		panic("query is empty")
	}

	if *config == "" {
		panic("config is empty")
	}

	reporter := mdq.DefaultReporter
	if *silent {
		reporter = mdq.SilentReporter
	}

	f, err := os.Open(*config)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dbs, err := mdq.CreateDBsFromConfig(f, *tag)
	if err != nil {
		panic(err)
	}

	cluster := mdq.NewCluster(dbs, reporter)

	results := cluster.Query(*query)

	var printer mdq.Printer
	if *format != "" {
		printer, err = mdq.NewTemplatePrinter(os.Stdout, *format)
		if err != nil {
			panic(err)
		}
	} else {
		printer = mdq.NewJsonPrinter(os.Stdout)
	}
	printer.Print(results)
}
