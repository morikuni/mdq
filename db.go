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

type DB interface {
	Query(query string) (Result, error)
}

func NewDB(name, driver, dsn string) (DB, error) {
	con, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return db{name, con}, nil
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
	defer rows.Close()

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
			v := values[i].Value
			row[col] = v
		}
		result.Rows = append(result.Rows, row)
	}
	err = rows.Err()
	if err != nil {
		return Result{}, db.err(err, "error in interating rows")
	}

	return result, nil
}

func (db db) err(err error, message string) error {
	return fmt.Errorf("%s: %s: %v", db.name, message, err)
}

type AnyValue struct {
	Value interface{}
}

func (v *AnyValue) Scan(src interface{}) error {
	v.Value = src
	if bs, ok := src.([]byte); ok {
		v.Value = string(bs)
	}
	return nil
}
