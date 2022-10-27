package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

// // ExecuteResult 执行结果
// type ExecuteResult interface {
// 	LastInsertId() (int64, error)
// 	RowsAffected() (int64, error)

// 	MapRecords2Map(rcverSupple func() []interface{},
// 		mapSupple func() map[string]interface{},
// 		mapFunc func(rcver []interface{}, dest map[string]interface{})) ([]map[string]interface{}, error)
// 	MapRecords2Struct(rcverSupple func() []interface{},
// 		mapFunc func(rcver []interface{}) interface{}) ([]interface{}, error)
// 	MapOneField(oneDest interface{}) error
// 	MapMultifield(multidest []interface{}) error

// 	Close() error
// 	HasError() error
// }

// type simpleExecuteResultImple struct {
// 	row       *sql.Row
// 	rows      *sql.Rows
// 	result    sql.Result
// 	execError error
// }

// func (imple *simpleExecuteResultImple) LastInsertId() (int64, error) {
// 	return imple.result.LastInsertId()
// }
// func (imple *simpleExecuteResultImple) RowsAffected() (int64, error) {
// 	return imple.result.RowsAffected()
// }

// func readMultirecord(rows *sql.Rows,
// 	rcver []interface{},
// 	handleRecord func(record []interface{})) error {

// 	for next := true; next; next = rows.NextResultSet() {
// 		for more := rows.Next(); more; more = rows.Next() {
// 			err := rows.Scan(rcver...)
// 			if err != nil {
// 				return err
// 			}
// 			handleRecord(rcver)
// 		}
// 	}

// 	return nil
// }
// func readOneRecord(row *sql.Row,
// 	rcver []interface{},
// 	handleRecord func(record []interface{})) error {

// 	err := row.Scan(rcver...)
// 	if err != nil {
// 		return err
// 	}
// 	handleRecord(rcver)

// 	return nil
// }
// func (imple *simpleExecuteResultImple) MapRecords2Map(rcverSupple func() []interface{},
// 	mapSupple func() map[string]interface{},
// 	mapFunc func(record []interface{}, dest map[string]interface{})) ([]map[string]interface{}, error) {
// 	ret := make([]map[string]interface{}, 0)
// 	if err := readMultirecord(imple.rows, rcverSupple(), func(record []interface{}) {
// 		dest := mapSupple()

// 		mapFunc(record, dest)
// 		ret = append(ret, dest)
// 	}); err != nil {
// 		return nil, err
// 	}

// 	return ret, nil
// }
// func (imple *simpleExecuteResultImple) MapRecords2Struct(rcverSupple func() []interface{},
// 	mapFunc func(rcver []interface{}) interface{}) ([]interface{}, error) {

// 	rcver := rcverSupple()
// 	rcverPtr := make([]interface{}, len(rcver))

// 	ret := make([]interface{}, 0)

// 	for next := true; next; next = imple.rows.NextResultSet() {
// 		for more := imple.rows.Next(); more; more = imple.rows.Next() {
// 			if err := imple.rows.Scan(rcverPtr...); err != nil {
// 				return nil, err
// 			}
// 			dest := mapFunc(rcver)
// 			ret = append(ret, dest)
// 		}
// 	}

// 	return ret, nil
// }
// func (imple *simpleExecuteResultImple) MapRecords2Struct2(rcverSupple func() []interface{},
// 	mapFunc func(rcver []interface{}) interface{}) ([]interface{}, error) {

// 	rcver := rcverSupple()
// 	rcverPtr := make([]interface{}, len(rcver))
// 	// settjing rcver 's ptr
// 	ret := make([]interface{}, 0)

// 	for next := true; next; next = imple.rows.NextResultSet() {
// 		for more := imple.rows.Next(); more; more = imple.rows.Next() {
// 			if err := imple.rows.Scan(rcverPtr...); err != nil {
// 				return nil, err
// 			}
// 			dest := mapFunc(rcver)
// 			ret = append(ret, dest)
// 		}
// 	}

// 	return ret, nil
// }
// func (imple *simpleExecuteResultImple) MapOneField(oneDest interface{}) error {
// 	rcver := []interface{}{oneDest}
// 	if err := readOneRecord(imple.row, rcver, func(record []interface{}) {

// 	}); err != nil {
// 		return err
// 	}
// 	return nil
// }
// func (imple *simpleExecuteResultImple) MapMultifield(multidest []interface{}) error {
// 	if err := readOneRecord(imple.row, multidest, func(record []interface{}) {

// 	}); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (imple *simpleExecuteResultImple) Close() error {
// 	if imple.rows == nil {
// 		return nil
// 	}
// 	return imple.rows.Close()
// }

// func (imple *simpleExecuteResultImple) HasError() error {
// 	return imple.execError
// }

// func newExecuteResult(oneRow *sql.Row, multirow *sql.Rows, execResult sql.Result, err error) ExecuteResult {
// 	return &simpleExecuteResultImple{
// 		row:       oneRow,
// 		rows:      multirow,
// 		result:    execResult,
// 		execError: err,
// 	}
// }

////////////////////////

// ExecuteResult sql 执行结果
type ExecuteResult struct {
	execError       error
	lastInsertID    int64
	lastInsertErr   error
	rowsAffected    int64
	rowsAffectedErr error
	row             *sql.Row
	rows            *sql.Rows
}

//HasExecuteError 是否有执行错误
func (result *ExecuteResult) HasExecuteError() error { return result.execError }

// LastInsertId 最新id
func (result *ExecuteResult) LastInsertId() (int64, error) {
	return result.lastInsertID, result.lastInsertErr
}

// RowsAffected 影响行数
func (result *ExecuteResult) RowsAffected() (int64, error) {
	return result.rowsAffected, result.rowsAffectedErr
}

// Close 关闭结果集数据流
func (result *ExecuteResult) Close() error {
	if result.rows == nil {
		return nil
	}
	return result.rows.Close()
}

// MapRecord2Struct 映射结果集到map
func (result *ExecuteResult) MapRecord2Struct(targetPtr interface{}) error {
	if targetPtr == nil {
		return errors.New("param must be non-nil pointer")
	}

	targetType := reflect.TypeOf(targetPtr)
	if targetType.Kind() != reflect.Ptr {
		return errors.New("param must be non-nil pointer")
	}
	targetRawType := targetType.Elem()

	fieldMappings := findFields(targetRawType, "*")

	//TODO resolve
	clmns, _ := result.rows.Columns()
	fmt.Printf("set clmns %v \n", clmns)
	rcver := make([]interface{}, len(clmns))
	for i, clmn := range clmns {
		rcver[i] = fieldMappings[clmn].FieldCache.Addr().Interface()
	}

	err := result.row.Scan(rcver...)
	if err != nil {
		return err
	}

	targetValue := reflect.ValueOf(targetPtr)
	rawValue := targetValue.Elem()
	for _, mapping := range fieldMappings {
		fieldValue := rawValue.FieldByName(mapping.StructField)
		fieldValue.Set(mapping.FieldCache)
	}
	return nil
}

// MapRecords2Struct 映射结果集到map
func (result *ExecuteResult) MapRecords2Struct(targetPtrSet interface{}) error {
	if targetPtrSet == nil {
		return errors.New("param must be non-nil pointer")
	}

	ptrType := reflect.TypeOf(targetPtrSet)
	if ptrType.Kind() != reflect.Ptr {
		return errors.New("param must be non-nil pointer")
	}
	slicType := ptrType.Elem()
	if slicType.Kind() != reflect.Slice {
		return errors.New("param must be a Slice")
	}

	eleType := slicType.Elem()
	var fieldMappings map[string]*FieldMapping
	if eleType.Kind() == reflect.Ptr {
		fieldMappings = findFields(eleType.Elem(), "*")
	} else {
		fieldMappings = findFields(eleType, "*")
	}

	var createInstance func() reflect.Value
	if eleType.Kind() == reflect.Ptr {
		createInstance = func() reflect.Value {
			return reflect.New(eleType.Elem()).Elem()
		}
	} else {
		createInstance = func() reflect.Value {
			return reflect.New(eleType).Elem()
		}
	}
	var convInstance func(reflect.Value) reflect.Value
	if eleType.Kind() == reflect.Ptr {
		convInstance = func(instance reflect.Value) reflect.Value {
			return instance.Addr()
		}
	} else {
		convInstance = func(instance reflect.Value) reflect.Value {
			return instance
		}
	}

	clmns, _ := result.rows.Columns()
	fmt.Printf("set clmns %v \n", clmns)
	rcver := make([]interface{}, len(clmns))
	for i, clmn := range clmns {
		rcver[i] = fieldMappings[clmn].FieldCache.Addr().Interface()
	}

	ptrValue := reflect.ValueOf(targetPtrSet)
	sliceValue := ptrValue.Elem()

	for more := result.rows.Next(); more; more = result.rows.Next() {
		if err := result.rows.Scan(rcver...); err != nil {
			return err
		}
		instanceValue := createInstance()
		for _, mapping := range fieldMappings {
			fieldValue := instanceValue.FieldByName(mapping.StructField)
			fieldValue.Set(mapping.FieldCache)
		}
		ele := convInstance(instanceValue)
		ns := reflect.Append(sliceValue, ele)
		sliceValue.Set(ns)
	}
	return nil
}

//MapRecord2Hashtable q
func (result *ExecuteResult) MapRecord2Hashtable(hashtbl map[string]interface{}) error {
	return nil
}
