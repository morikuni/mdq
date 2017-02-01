package mdq

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDB struct {
	Result Result
	Err    error
}

func (db testDB) Query(query string) (Result, error) {
	return db.Result, db.Err
}

func CreateTestDB(r Result, err error) DB {
	return testDB{r, err}
}

type testReporter struct {
	Err error
}

func (r *testReporter) Report(err error) {
	r.Err = err
}

func ToMap(results []Result) map[string]Result {
	m := make(map[string]Result)
	for _, r := range results {
		m[r.Database] = r
	}
	return m
}

func TestCluster(t *testing.T) {

	type Input struct {
		DBs []DB
	}
	type Output struct {
		Results     []Result
		ReportedErr error
	}
	type Test struct {
		Title  string
		Input  Input
		Output Output
	}

	table := []Test{
		Test{
			Title: "success",
			Input: Input{
				DBs: []DB{
					CreateTestDB(
						Result{
							Database: "db1",
							Columns:  []string{"id"},
							Rows: []map[string]interface{}{
								map[string]interface{}{
									"id": 1,
								},
							},
						},
						nil,
					),
					CreateTestDB(
						Result{
							Database: "db2",
							Columns:  []string{"name"},
							Rows: []map[string]interface{}{
								map[string]interface{}{
									"name": "foo",
								},
							},
						},
						nil,
					),
				},
			},
			Output: Output{
				Results: []Result{
					Result{
						Database: "db1",
						Columns:  []string{"id"},
						Rows: []map[string]interface{}{
							map[string]interface{}{
								"id": 1,
							},
						},
					},
					Result{
						Database: "db2",
						Columns:  []string{"name"},
						Rows: []map[string]interface{}{
							map[string]interface{}{
								"name": "foo",
							},
						},
					},
				},
				ReportedErr: nil,
			},
		},
		Test{
			Title: "report error when some db returns error",
			Input: Input{
				DBs: []DB{
					CreateTestDB(
						Result{
							Database: "db1",
							Columns:  []string{"id"},
							Rows: []map[string]interface{}{
								map[string]interface{}{
									"id": 1,
								},
							},
						},
						nil,
					),
					CreateTestDB(Result{}, errors.New("query error")),
				},
			},
			Output: Output{
				Results: []Result{
					Result{
						Database: "db1",
						Columns:  []string{"id"},
						Rows: []map[string]interface{}{
							map[string]interface{}{
								"id": 1,
							},
						},
					},
				},
				ReportedErr: errors.New("query error"),
			},
		},
	}

	for _, test := range table {
		t.Run(test.Title, func(t *testing.T) {
			assert := assert.New(t)

			reporter := &testReporter{}

			cluster := NewCluster(test.Input.DBs, reporter)

			r := cluster.Query("query")

			// Use ToMap to fix slice order problem
			assert.Equal(ToMap(test.Output.Results), ToMap(r))
			assert.Equal(test.Output.ReportedErr, reporter.Err)
		})
	}
}
