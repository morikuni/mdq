package mdq

import (
	"encoding/json"
	"html/template"
	"io"
)

type Printer interface {
	Print(w io.Writer, results []Result) error
}

func NewTemplatePrinter(format string) (Printer, error) {
	t, err := template.New("sql").Parse(format)
	if err != nil {
		return nil, err
	}
	return templatePrinter{t}, nil
}

type templatePrinter struct {
	t *template.Template
}

func (p templatePrinter) Print(w io.Writer, results []Result) error {
	return p.t.Execute(w, results)
}

func NewJsonPrinter() Printer {
	return jsonPrinter{}
}

type jsonPrinter struct{}

func (p jsonPrinter) Print(w io.Writer, results []Result) error {
	return json.NewEncoder(w).Encode(results)
}
