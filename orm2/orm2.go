package orm2

import (
	"database/sql"
	"errors"
	"reflect"
)

// ExecuteResult2 sql 执行结果
type ExecuteResult2 struct {
	execError       error
	lastInsertID    int64
	lastInsertErr   error
	rowsAffected    int64
	rowsAffectedErr error
	row             *sql.Row
	rows            *sql.Rows
}

//HasExecuteError 是否有执行错误
func (result *ExecuteResult2) HasExecuteError() error { return result.execError }

// LastInsertId 最新id
func (result *ExecuteResult2) LastInsertId() (int64, error) {
	return result.lastInsertID, result.lastInsertErr
}

// RowsAffected 影响行数
func (result *ExecuteResult2) RowsAffected() (int64, error) {
	return result.rowsAffected, result.rowsAffectedErr
}

// Close 关闭结果集数据流
func (result *ExecuteResult2) Close() error {
	if result.rows == nil {
		return nil
	}
	return result.rows.Close()
}

// MapRecord2Struct 映射结果集到map
func (result *ExecuteResult2) MapRecord2Struct(targetPtr interface{}) error {
	if targetPtr == nil {
		return errors.New("param must be non-nil pointer")
	}

	targetType := reflect.TypeOf(targetPtr)
	if targetType.Kind() != reflect.Ptr {
		return errors.New("param must be non-nil pointer")
	}
	targetRawType := targetType.Elem()

	fieldMappings := findFields(targetRawType, "*")
	rcver := make([]interface{}, len(fieldMappings))
	for i, mapping := range fieldMappings {
		rcver[i] = mapping.FieldCache.Addr().Interface()
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
func (result *ExecuteResult2) MapRecords2Struct(targetPtrSet interface{}) error {
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
	var fieldMappings []*FieldMapping
	if eleType.Kind() == reflect.Ptr {
		fieldMappings = findFields(eleType.Elem(), "*")
	} else {
		fieldMappings = findFields(eleType, "*")
	}

	var createInstance func() reflect.Value
	if eleType.Kind() == reflect.Ptr {
		createInstance = func() reflect.Value {
			return reflect.New(eleType.Elem())
		}
	} else {
		createInstance = func() reflect.Value {
			return reflect.New(eleType)
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

	rcver := make([]interface{}, len(fieldMappings))
	for i, mapping := range fieldMappings {
		rcver[i] = mapping.FieldCache.Addr().Interface()
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
func (result *ExecuteResult2) MapRecord2Hashtable(hashtbl map[string]interface{}) error {
	return nil
}

type Executer interface {
	WriteRecords(sql string) ExecuteResult2
	ReadOneRecord(sql string) ExecuteResult2
	ReadRecords(sql string) ExecuteResult2
}
type Query interface {
	From(tbl string) Query
	Select(fieldName ...string) Query
	Count() Query
	Sum(fieldName string) Query
	Max(fieldName string) Query
	Min(fieldName string) Query
	Avg(fieldName string) Query
	Where() Query
	Group(fieldName ...string) Query
	Order() Query
	Offset(o int) Query
	Limit(l int) Query
	MapTo(target interface{}) Query
	Error() error
}
