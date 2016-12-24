package mdq

import (
	"database/sql"
	"fmt"
)

type Executor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Result struct {
	Database string
	Columns  []string
	Rows     []map[string]interface{}
}

type column struct {
	value interface{}
}

func (c *column) Scan(src interface{}) error {
	c.value = src
	if bs, ok := src.([]byte); ok {
		c.value = string(bs)
	}
	return nil
}

type DB interface {
	Query(query string) (Result, error)
}

func NewDB(name string, executor Executor) DB {
	return db{name, executor}
}

type db struct {
	name     string
	executor Executor
}

func (db db) Query(query string) (Result, error) {
	rows, err := db.executor.Query(query)
	if err != nil {
		return Result{}, db.err(err, "execution failed")
	}

	columns, err := rows.Columns()
	if err != nil {
		return Result{}, db.err(err, "failed to fetch column name")
	}
	columnSize := len(columns)

	result := Result{
		Database: db.name,
		Columns:  columns,
	}
	for rows.Next() {
		values := make([]AnyValue, columnSize)
		valuePtrs := make([]interface{}, columnSize)
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return Result{}, db.err(err, "failed to bind values")
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			v := values[i].val
			row[col] = v
		}
		result.Rows = append(result.Rows, row)
	}
	return result, nil
}

func (db db) err(err error, message string) error {
	return fmt.Errorf("%s: %s: %v", db.name, message, err)
}

type AnyValue struct {
	val interface{}
}

func (v *AnyValue) Scan(src interface{}) error {
	v.val = src
	if bs, ok := src.([]byte); ok {
		v.val = string(bs)
	}
	return nil
}
