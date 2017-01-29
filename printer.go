package mdq

import (
	"encoding/json"
	"html/template"
	"io"
)

type Printer interface {
	Print(results []Result) error
}

func NewTemplatePrinter(w io.Writer, format string) (Printer, error) {
	t, err := template.New("sql").Parse(format)
	if err != nil {
		return nil, err
	}
	return templatePrinter{w, t}, nil
}

type templatePrinter struct {
	w io.Writer
	t *template.Template
}

func (p templatePrinter) Print(results []Result) error {
	return p.t.Execute(p.w, results)
}

func NewJsonPrinter(w io.Writer) Printer {
	return jsonPrinter{
		json.NewEncoder(w),
	}
}

type jsonPrinter struct {
	encoder *json.Encoder
}

func (p jsonPrinter) Print(results []Result) error {
	return p.encoder.Encode(results)
}
