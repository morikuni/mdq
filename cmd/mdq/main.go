package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/morikuni/mdq"
	"github.com/spf13/pflag"
)

var (
	Version string = "unknown"
)

func main() {
	os.Exit(Run(os.Args, os.Stdin, os.Stdout, os.Stderr))
}

func Run(args []string, in io.Reader, out io.Writer, errW io.Writer) int {
	home := os.Getenv("HOME")

	flag := pflag.NewFlagSet("mdq", pflag.ContinueOnError)
	tag := flag.String("tag", "", "database tag")
	format := flag.String("format", "", "golang template string")
	query := flag.StringP("query", "q", "", "SQL")
	config := flag.String("config", home+"/.config/mdq/config.yaml", "path to config file")
	silent := flag.Bool("silent", false, "ignore errors from databases")
	help := flag.BoolP("help", "h", false, "print this help.")
	version := flag.Bool("version", false, "print version of mdq")

	flag.Usage = func() {
		fmt.Fprintln(errW)
		fmt.Fprintln(errW, "Usage: mdq [flags]")
		fmt.Fprintln(errW)
		fmt.Fprintln(errW, flag.FlagUsages())
	}

	err := flag.Parse(args[1:])
	if err != nil {
		fmt.Fprintln(errW, err)
		return 1
	}

	if *help {
		flag.Usage()
		return 0
	}

	if *version {
		fmt.Fprintln(out, "mdq version", Version)
		return 0
	}

	if *query == "" {
		bs, err := ioutil.ReadAll(in)
		*query = string(bs)
		if err != nil {
			fmt.Fprintln(errW, *config)
			return 1
		}
	}

	reporter := mdq.DefaultReporter
	if *silent {
		reporter = mdq.SilentReporter
	}

	f, err := os.Open(*config)
	if err != nil {
		fmt.Fprintln(errW, "cannot open config file:", *config)
		return 1
	}
	defer f.Close()

	dbs, err := mdq.CreateDBsFromConfig(f, *tag)
	if err != nil {
		fmt.Fprintln(errW, err)
		return 1
	}

	cluster := mdq.NewCluster(dbs, reporter)

	results := cluster.Query(*query)

	var printer mdq.Printer
	if *format != "" {
		printer, err = mdq.NewTemplatePrinter(os.Stdout, *format)
		if err != nil {
			fmt.Fprintln(errW, err)
			return 1
		}
	} else {
		printer = mdq.NewJsonPrinter(os.Stdout)
	}
	printer.Print(results)

	return 0
}
