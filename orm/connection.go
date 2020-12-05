package orm

import (
	"database/sql"
	"fmt"
	"strings"
)

// DatabaseConnection 连接
type DatabaseConnection interface {
	IsConnecting() bool

	ExecuteCMD(sql string) ExecuteResult

	CreateRecords(fields []string, table, where string, values [][]interface{}) ExecuteResult
	DeleteRecords(table, where string, parameters []interface{}) ExecuteResult
	UpdateRecords(updates []string, values []interface{}, table, where string, parameters []interface{}) ExecuteResult
	QueryMultirecord(fields []string, table, where string, parameters []interface{}) ExecuteResult
	QueryOneRecord(fields []string, table, where string, parameters []interface{}) ExecuteResult
	QueryMultitable(fields []string, mainTable, where string, parameters []interface{}, pointers []*MultitableJoinpointer) ExecuteResult
}
type simpleDBConnImple struct {
	conn *sql.DB
}

func (impl *simpleDBConnImple) IsConnecting() bool {
	return impl.conn.Ping() == nil
}

func (impl *simpleDBConnImple) ExecuteCMD(sql string) ExecuteResult {
	rst, err := impl.conn.Exec(sql)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, nil, rst, nil)
}

func (impl *simpleDBConnImple) CreateRecords(fields []string, table, where string, values [][]interface{}) ExecuteResult {

	if len(values) == 0 {
		return newExecuteResult(nil, nil, nil, ERR_CREATE_PARAM_MUST_NOT_BE_EMPTY)
	}
	if len(values[0]) != len(fields) {
		return newExecuteResult(nil, nil, nil, ERR_CREATE_LEN_FIELDS_MUST_SAME_VALUES)
	}

	sqlFragments := make([]string, len(values)+1)

	fills := strings.Join(fields, ",")
	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES", table, fills)
	sqlFragments[0] = insertSQL

	valuesTemplate := "(%s)"
	for i, record := range values {
		value, err := stringValues(record)
		if err != nil {
			return newExecuteResult(nil, nil, nil, err)
		}
		fragment := fmt.Sprintf(valuesTemplate, value)
		sqlFragments[i+1] = fragment
	}

	sql := strings.Join(sqlFragments, ",")
	rst, err := impl.conn.Exec(sql)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, nil, rst, nil)
}
func (impl *simpleDBConnImple) DeleteRecords(table, where string, parameters []interface{}) ExecuteResult {

	for _, p := range parameters {
		s, e := stringValue(p)
		if e != nil {
			return newExecuteResult(nil, nil, nil, e)
		}
		where = strings.Replace(where, "?", s, 0)
	}
	deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)

	rst, err := impl.conn.Exec(deleteSQL)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, nil, rst, nil)
}
func (impl *simpleDBConnImple) UpdateRecords(updates []string, values []interface{}, table, where string, parameters []interface{}) ExecuteResult {
	if len(values) == 0 {
		return newExecuteResult(nil, nil, nil, ERR_CREATE_PARAM_MUST_NOT_BE_EMPTY)
	}
	if len(values) != len(updates) {
		return newExecuteResult(nil, nil, nil, ERR_CREATE_LEN_FIELDS_MUST_SAME_VALUES)
	}

	for _, p := range parameters {
		s, e := stringValue(p)
		if e != nil {
			return newExecuteResult(nil, nil, nil, e)
		}
		where = strings.Replace(where, "?", s, 0)
	}

	updateTemplate := "%s = %s"

	updateFragments := make([]string, len(updates))
	for i := 0; i < len(updates); i++ {

		value, err := stringValue(values[i])
		if err != nil {
			return newExecuteResult(nil, nil, nil, err)
		}
		fragment := fmt.Sprintf(updateTemplate, updates[i], value)
		updateFragments[i] = fragment
	}

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", strings.Join(updateFragments, ","), table, where)
	rst, err := impl.conn.Exec(sql)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, nil, rst, nil)
}
func (impl *simpleDBConnImple) QueryMultirecord(selects []string, table, where string, parameters []interface{}) ExecuteResult {
	if len(selects) == 0 {
		return newExecuteResult(nil, nil, nil, ERR_QUERY_SELECT_CANOT_BE_BLANK)
	}

	for _, p := range parameters {
		s, e := stringValue(p)
		if e != nil {
			return newExecuteResult(nil, nil, nil, e)
		}
		where = strings.Replace(where, "?", s, 0)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(selects, ","), table, where)
	rows, err := impl.conn.Query(sql)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, rows, nil, nil)
}
func (impl *simpleDBConnImple) QueryOneRecord(selects []string, table, where string, parameters []interface{}) ExecuteResult {
	if len(selects) == 0 {
		return newExecuteResult(nil, nil, nil, ERR_QUERY_SELECT_CANOT_BE_BLANK)
	}

	for _, p := range parameters {
		s, e := stringValue(p)
		if e != nil {
			return newExecuteResult(nil, nil, nil, e)
		}
		where = strings.Replace(where, "?", s, 0)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(selects, ","), table, where)
	row := impl.conn.QueryRow(sql)

	return newExecuteResult(row, nil, nil, nil)
}

//MultitableJoinpointer 多表联接
type MultitableJoinpointer struct {
	Join    string
	Table   string
	JoinExp string
}

func (jp *MultitableJoinpointer) String() string {
	return fmt.Sprintf("%s JOIN %s ON  %s ", jp.Join, jp.Table, jp.JoinExp)
}

func (impl *simpleDBConnImple) QueryMultitable(selects []string, mainTable, where string, parameters []interface{}, pointers []*MultitableJoinpointer) ExecuteResult {
	if len(selects) == 0 {
		return newExecuteResult(nil, nil, nil, ERR_QUERY_SELECT_CANOT_BE_BLANK)
	}

	joinTables := make([]string, len(pointers))
	for i, j := range pointers {
		joinTables[i] = j.String()
	}

	for _, p := range parameters {
		s, e := stringValue(p)
		if e != nil {
			return newExecuteResult(nil, nil, nil, e)
		}
		where = strings.Replace(where, "?", s, 0)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s %s  WHERE %s", strings.Join(selects, ","), mainTable, strings.Join(joinTables, " "), where)
	rows, err := impl.conn.Query(sql)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}
	return newExecuteResult(nil, rows, nil, nil)
}

// ConnectionProperties 链接属性
type ConnectionProperties struct {
	Dialect        string
	User           string
	Password       string
	ConnectionType string
	Host           string
	Port           uint16
	Database       string
	Charset        string
}

//OpenConnection 打开
func OpenConnection(prop ConnectionProperties) (DatabaseConnection, error) {
	dialect := prop.Dialect
	url := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s",
		prop.User, prop.Password,
		prop.ConnectionType, prop.Host, prop.Port,
		prop.Database, prop.Charset)
	db, err := sql.Open(dialect, url)
	if err != nil {
		return nil, err
	}
	return &simpleDBConnImple{conn: db}, nil
}

// 偏移，限制
// 分组
// 排序
// union
// 值和类型反射
// 事务 aop

// Storeable 数据反射
type Storeable interface {
}

// Tx 事务契约
type Tx interface {
}
