package mdq

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {

	type Input struct {
		QueryErr error
		Columns  []string
		Values   [][]driver.Value
	}
	type Output struct {
		Result Result
		Err    error
	}
	type Test struct {
		Title  string
		Input  Input
		Output Output
	}

	name := "test_database"
	table := []Test{
		Test{
			Title: "success",
			Input: Input{
				QueryErr: nil,
				Columns:  []string{"id", "name"},
				Values: [][]driver.Value{
					[]driver.Value{1, "foo"},
					[]driver.Value{2, "bar"},
				},
			},
			Output: Output{
				Result: Result{
					Database: name,
					Columns:  []string{"id", "name"},
					Rows: []map[string]interface{}{
						map[string]interface{}{
							"id":   1,
							"name": "foo",
						},
						map[string]interface{}{
							"id":   2,
							"name": "bar",
						},
					},
				},
				Err: nil,
			},
		},
		Test{
			Title: "error when Query returns error",
			Input: Input{
				QueryErr: errors.New("query error"),
				Columns:  nil,
				Values:   nil,
			},
			Output: Output{
				Result: Result{},
				Err:    fmt.Errorf("%s: %s: %s", name, "execution failed", "query error"),
			},
		},
	}

	for _, test := range table {
		t.Run(test.Title, func(t *testing.T) {
			assert := assert.New(t)

			con, mock, err := sqlmock.New()
			assert.NoError(err)

			rows := sqlmock.NewRows(test.Input.Columns)
			for _, values := range test.Input.Values {
				rows.AddRow(values...)
			}

			mock.ExpectQuery("query").WillReturnRows(rows).WillReturnError(test.Input.QueryErr)

			db := db{name, con}

			r, err := db.Query("query")

			assert.Equal(test.Output.Err, err)
			assert.Equal(test.Output.Result, r)
		})
	}
}

func TestAnyValue(t *testing.T) {

	type Input struct {
		Value interface{}
	}
	type Output struct {
		Value interface{}
	}
	type Test struct {
		Title  string
		Input  Input
		Output Output
	}

	now := time.Now()
	table := []Test{
		Test{
			Title:  "int64",
			Input:  Input{int64(12345)},
			Output: Output{int64(12345)},
		},
		Test{
			Title:  "float64",
			Input:  Input{float64(12345)},
			Output: Output{float64(12345)},
		},
		Test{
			Title:  "bool",
			Input:  Input{true},
			Output: Output{true},
		},
		Test{
			Title:  "[]byte",
			Input:  Input{[]byte("hello")},
			Output: Output{"hello"},
		},
		Test{
			Title:  "string",
			Input:  Input{"hello"},
			Output: Output{"hello"},
		},
		Test{
			Title:  "time.Time",
			Input:  Input{now},
			Output: Output{now},
		},
		Test{
			Title:  "nil",
			Input:  Input{nil},
			Output: Output{nil},
		},
	}

	for _, test := range table {
		t.Run(test.Title, func(t *testing.T) {
			assert := assert.New(t)

			var a AnyValue

			err := a.Scan(test.Input.Value)

			assert.NoError(err)
			assert.Equal(test.Output.Value, a.Value)
		})
	}
}
