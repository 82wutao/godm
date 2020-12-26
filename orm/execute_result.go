package orm

import (
	"database/sql"
)

// ExecuteResult 执行结果
type ExecuteResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)

	MapRecords2Map(rcverSupple func() []interface{},
		mapSupple func() map[string]interface{},
		mapFunc func(rcver []interface{}, dest map[string]interface{})) ([]map[string]interface{}, error)
	MapRecords2Struct(rcverSupple func() []interface{},
		structSupple func() interface{},
		mapFunc func(rcver []interface{}, dest interface{})) ([]interface{}, error)
	MapOneField(oneDest interface{}) error
	MapMultifield(multidest []interface{}) error

	Close() error
	HasError() error
}

type simpleExecuteResultImple struct {
	row       *sql.Row
	rows      *sql.Rows
	result    sql.Result
	execError error
}

func (imple *simpleExecuteResultImple) LastInsertId() (int64, error) {
	return imple.result.LastInsertId()
}
func (imple *simpleExecuteResultImple) RowsAffected() (int64, error) {
	return imple.result.RowsAffected()
}

func readMultirecord(rows *sql.Rows,
	rcver []interface{},
	handleRecord func(record []interface{})) error {

	for next := true; next; next = rows.NextResultSet() {
		for more := rows.Next(); more; more = rows.Next() {
			err := rows.Scan(rcver...)
			if err != nil {
				return err
			}
			handleRecord(rcver)
		}
	}

	return nil
}
func readOneRecord(row *sql.Row,
	rcver []interface{},
	handleRecord func(record []interface{})) error {

	err := row.Scan(rcver...)
	if err != nil {
		return err
	}
	handleRecord(rcver)

	return nil
}
func (imple *simpleExecuteResultImple) MapRecords2Map(rcverSupple func() []interface{},
	mapSupple func() map[string]interface{},
	mapFunc func(record []interface{}, dest map[string]interface{})) ([]map[string]interface{}, error) {
	ret := make([]map[string]interface{}, 0)
	if err := readMultirecord(imple.rows, rcverSupple(), func(record []interface{}) {
		dest := mapSupple()

		mapFunc(record, dest)
		ret = append(ret, dest)
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
func (imple *simpleExecuteResultImple) MapRecords2Struct(rcverSupple func() []interface{},
	structSupple func() interface{},
	mapFunc func(record []interface{}, dest interface{})) ([]interface{}, error) {

	ret := make([]interface{}, 0)
	if err := readMultirecord(imple.rows, rcverSupple(), func(record []interface{}) {
		dest := structSupple()

		mapFunc(record, dest)
		ret = append(ret, dest)
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
func (imple *simpleExecuteResultImple) MapOneField(oneDest interface{}) error {
	rcver := []interface{}{oneDest}
	if err := readOneRecord(imple.row, rcver, func(record []interface{}) {

	}); err != nil {
		return err
	}
	return nil
}
func (imple *simpleExecuteResultImple) MapMultifield(multidest []interface{}) error {
	if err := readOneRecord(imple.row, multidest, func(record []interface{}) {

	}); err != nil {
		return err
	}
	return nil
}

func (imple *simpleExecuteResultImple) Close() error {
	if imple.rows == nil {
		return nil
	}
	return imple.rows.Close()
}

func (imple *simpleExecuteResultImple) HasError() error {
	return imple.execError
}

func newExecuteResult(oneRow *sql.Row, multirow *sql.Rows, execResult sql.Result, err error) ExecuteResult {
	return &simpleExecuteResultImple{
		row:       oneRow,
		rows:      multirow,
		result:    execResult,
		execError: err,
	}
}
