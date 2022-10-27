package orm

import "dm.net/datamine/orm/clause"

type Executer interface {
	WriteRecords(sql string) ExecuteResult
	ReadOneRecord(sql string) ExecuteResult
	ReadRecords(sql string) ExecuteResult
}

// Save save it to base
type Save interface {
	// Save  {struct,slic}
	Save(ptr interface{}) ([]int64, error)
}

type Detele interface {
	From(tbl string) Detele
	Delete(where clause.WhereClause) (int64, error)
}
type Update interface {
	From(tbl string) Update
	Select(fieldName ...string) Update
	Count() Update
	Sum(fieldName string) Update
	Max(fieldName string) Update
	Min(fieldName string) Update
	Avg(fieldName string) Update
	Where() Update
	Group(fieldName ...string) Update
	Order() Update
	Offset(o int) Update
	Limit(l int) Update
	MapTo(target interface{}) Update
	Error() error
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

type JoinQuery interface {
	From(tbl string) JoinQuery
	Select(fieldName ...string) JoinQuery
	Count() JoinQuery
	Sum(fieldName string) JoinQuery
	Max(fieldName string) JoinQuery
	Min(fieldName string) JoinQuery
	Avg(fieldName string) JoinQuery
	Where() JoinQuery
	Group(fieldName ...string) JoinQuery
	Order() JoinQuery
	Offset(o int) JoinQuery
	Limit(l int) JoinQuery
	MapTo(target interface{}) JoinQuery
	Error() error
}
