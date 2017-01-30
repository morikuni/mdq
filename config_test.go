package mdq

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyDB struct{}

func (db dummyDB) Query(q string) (Result, error) {
	return Result{}, nil
}

func TestCreateDBsFromFile(t *testing.T) {

	type Input struct {
		Text     string
		Tag      string
		NewDBErr error
	}
	type Output struct {
		DBsLen int
		Err    error
	}
	type Test struct {
		Title  string
		Input  Input
		Output Output
	}

	yaml := `dbs:
  - name: "db1"
    driver: "mysql"
    dsn: "root@/db"
    tags: ["hello"]
  - name: "db2"
    driver: "sqlite3"
    dsn: "/path/to/db.sqlite3"
    tags: ["world"]`

	table := []Test{
		Test{
			Title: "success with no tag",
			Input: Input{
				Text:     yaml,
				Tag:      "",
				NewDBErr: nil,
			},
			Output: Output{
				DBsLen: 2,
				Err:    nil,
			},
		},
		Test{
			Title: "success with specifying tag",
			Input: Input{
				Text:     yaml,
				Tag:      "hello",
				NewDBErr: nil,
			},
			Output: Output{
				DBsLen: 1,
				Err:    nil,
			},
		},
		Test{
			Title: "error when NewDB returns error",
			Input: Input{
				Text:     yaml,
				Tag:      "",
				NewDBErr: errors.New("errrr"),
			},
			Output: Output{
				DBsLen: 0,
				Err:    errors.New("errrr"),
			},
		},
	}

	for _, test := range table {
		t.Run(test.Title, func(t *testing.T) {
			assert := assert.New(t)
			cp := ConfigParser{
				func(_, _, _ string) (DB, error) {
					return dummyDB{}, test.Input.NewDBErr
				},
			}
			dbs, err := cp.CreateDBsFromConfig(strings.NewReader(test.Input.Text), test.Input.Tag)
			assert.Len(dbs, test.Output.DBsLen)
			assert.Equal(test.Output.Err, err)
		})
	}
}
