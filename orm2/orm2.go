package orm2

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
