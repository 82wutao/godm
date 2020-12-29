package orm

import (
	"database/sql"
	"fmt"
	"strings"

	"dm.net/datamine/orm/clause"
	"dm.net/datamine/orm/errors"
	"dm.net/datamine/orm/util"
	"dm.net/datamine/syntaxutil"
)

// DatabaseConnection 连接
type DatabaseConnection interface {
	IsConnecting() bool

	ExecuteCMD(sql string) ExecuteResult

	CreateRecords(mapping StructToRecordMapping, structs []interface{}) ExecuteResult
	DeleteRecords(table string, where clause.WhereClause) ExecuteResult
	UpdateRecords(mapping StructToRecordMapping, data interface{}, where clause.WhereClause) ExecuteResult
	QueryMultirecord(mapping RecordToStructMapping) ExecuteResult
	QueryOneRecord(mapping RecordToStructMapping) ExecuteResult
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

func (impl *simpleDBConnImple) CreateRecords(mapping StructToRecordMapping, structs []interface{}) ExecuteResult {

	if len(structs) == 0 {
		return newExecuteResult(nil, nil, nil, errors.ERR_CREATE_PARAM_MUST_NOT_BE_EMPTY)
	}

	fragments := make([]string, len(structs))
	for i, s := range structs {
		values := mapping.MapsStructToValues(s)
		fragment := util.StringValues(values)
		fragments[i] = fmt.Sprintf("(%s)", fragment)
	}
	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		mapping.DataSourceMapped(), strings.Join(mapping.FieldsMapped(), ","), strings.Join(fragments, ","))

	rst, err := impl.conn.Exec(insertSQL)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, nil, rst, nil)
}
func (impl *simpleDBConnImple) DeleteRecords(table string, where clause.WhereClause) ExecuteResult {
	whereExp := syntaxutil.TernaryOperate(where == nil, "", func() interface{} { return where.WhereSQL() }).(string)

	deleteSQL := fmt.Sprintf("DELETE FROM %s %s", table, whereExp)
	rst, err := impl.conn.Exec(deleteSQL)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, nil, rst, nil)
}
func (impl *simpleDBConnImple) UpdateRecords(mapping StructToRecordMapping, data interface{}, where clause.WhereClause) ExecuteResult {

	clmns := mapping.FieldsMapped()
	updates := make([]string, len(clmns))

	for i, p := range mapping.MapsStructToValues(data) {
		v := util.StringValue(p)

		updates[i] = fmt.Sprintf("%s = %s", clmns[i], v)
	}

	whereExp := syntaxutil.TernaryOperate(where == nil, "", func() interface{} { return where.WhereSQL() })
	sql := fmt.Sprintf("UPDATE %s SET %s %s", mapping.DataSourceMapped(), strings.Join(updates, ","), whereExp)
	rst, err := impl.conn.Exec(sql)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, nil, rst, nil)
}

func (impl *simpleDBConnImple) QueryMultirecord(mapping RecordToStructMapping) ExecuteResult {

	whereExp := syntaxutil.TernaryOperate(mapping.WhereClause == nil, "", func() interface{} { return mapping.WhereClause.WhereSQL() })
	groupExp := syntaxutil.TernaryOperate(mapping.GroupClause == nil, "", func() interface{} { return mapping.GroupClause.GroupSQL() })
	orderExp := syntaxutil.TernaryOperate(mapping.OrderClause == nil, "", func() interface{} { return mapping.OrderClause.OrderSQL() })
	offsetLimitExp := syntaxutil.TernaryOperate(mapping.OffsetLimitClause == nil, "", func() interface{} { return mapping.OffsetLimitClause.OffsetLimitSQL() })
	selectSQL := fmt.Sprintf("SELECT %s FROM %s %s %s %s %s",
		strings.Join(mapping.FieldsMapped(), ","), mapping.DataSourceMapped(),
		whereExp, groupExp, orderExp, offsetLimitExp)
	rows, err := impl.conn.Query(selectSQL)
	if err != nil {
		return newExecuteResult(nil, nil, nil, err)
	}

	return newExecuteResult(nil, rows, nil, nil)
}
func (impl *simpleDBConnImple) QueryOneRecord(mapping RecordToStructMapping) ExecuteResult {
	whereExp := syntaxutil.TernaryOperate(mapping.WhereClause == nil, "", func() interface{} { return mapping.WhereClause.WhereSQL() })
	groupExp := syntaxutil.TernaryOperate(mapping.GroupClause == nil, "", func() interface{} { return mapping.GroupClause.GroupSQL() })
	selectSQL := fmt.Sprintf("SELECT %s FROM %s %s %s",
		strings.Join(mapping.FieldsMapped(), ","), mapping.DataSourceMapped(), whereExp, groupExp)
	row := impl.conn.QueryRow(selectSQL)

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
		return newExecuteResult(nil, nil, nil, errors.ERR_QUERY_SELECT_CANOT_BE_BLANK)
	}

	joinTables := make([]string, len(pointers))
	for i, j := range pointers {
		joinTables[i] = j.String()
	}

	for _, p := range parameters {
		s := util.StringValue(p)
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
	// connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	// db, err := sql.Open("postgres", connStr)

	dialect := prop.Dialect
	url := fmt.Sprintf("%s://%s:%s@%s:%d/%s?charset=%s&sslmode=disable",
		prop.Dialect,
		prop.User, prop.Password,
		prop.Host, prop.Port,
		prop.Database, prop.Charset)
	db, err := sql.Open(dialect, url)
	if err != nil {
		return nil, err
	}
	return &simpleDBConnImple{conn: db}, nil
}
