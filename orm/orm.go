package orm

import "dm.net/datamine/orm/clause"

//StructToRecordMapping 可以被记录的对象 ,用在create 和 update
type StructToRecordMapping interface {
	FieldsMapped() clause.FieldsMapped
	DataSourceMapped() clause.DataSourceMapped
	MapsStructToValues(strct interface{}) []interface{}
}

//RecordToStructMapping 数据库记录到内存结构
type RecordToStructMapping interface {
	FieldsMapped() clause.FieldsMapped
	DataSourceMapped() clause.DataSourceMapped
	WhereClause() clause.WhereClause
	GroupClause() clause.GroupClause
	OrderClause() clause.OrderClause
	OffsetLimitClause() clause.OffsetLimitClause
	Collector() []interface{}
	Map(collector []interface{}) interface{}
}
