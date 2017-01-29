package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/morikuni/mdq"
)

func main() {
	home := os.Getenv("HOME")

	targetReg := kingpin.Flag("target", "reqular expression to filter target databases").Regexp()
	format := kingpin.Flag("format", "output format").Short('f').String()
	query := kingpin.Flag("query", "SQL").Short('q').String()
	config := kingpin.Flag("config", "path to config file").Short('c').Default(home + "/.config/mdq/config.yaml").String()
	silent := kingpin.Flag("silent", "ignore errors from databases").Short('s').Default("false").Bool()
	kingpin.Parse()

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

	dbs, err := mdq.CreateDBsFromFile(f, *targetReg)
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
