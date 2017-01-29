package main

import (
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/morikuni/mdq"
	"github.com/spf13/pflag"
)

func main() {
	home := os.Getenv("HOME")

	flag := pflag.NewFlagSet("mdq", pflag.ContinueOnError)
	target := flag.String("target", "", "target filtering regular expression")
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

	var targetReg *regexp.Regexp
	if *target != "" {
		targetReg, err = regexp.Compile(*target)
		if err != nil {
			panic(err)
		}
	}
	dbs, err := mdq.CreateDBsFromFile(f, targetReg)
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
