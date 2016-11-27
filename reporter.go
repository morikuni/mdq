package mdq

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	DefaultReporter = NewReporter(os.Stderr)
	SilentReporter  = NewReporter(ioutil.Discard)
)

type Reporter interface {
	Report(err error)
}

func NewReporter(w io.Writer) Reporter {
	return reporter{w}
}

type reporter struct {
	w io.Writer
}

func (r reporter) Report(err error) {
	fmt.Fprintln(r.w, err)
}
