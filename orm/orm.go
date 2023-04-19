package orm

import (
	"dm.net/datamine/common"
	"dm.net/datamine/orm/clause"
)

//StructToRecordMapping 可以被记录的对象 ,用在create 和 update
type StructToRecordMapping struct {
	FieldsMapped       clause.TargetFields
	DataSourceMapped   clause.TargetDataSource
	MapsStructToValues common.Func[common.T, []common.R]
}

type RecordToStructMapping struct {
	FieldsMapped      clause.TargetFields
	DataSourceMapped  clause.TargetDataSource
	WhereClause       clause.WhereClause
	GroupClause       clause.GroupClause
	OrderClause       clause.OrderClause
	OffsetLimitClause clause.OffsetLimitClause
	Collector         func() []interface{}
	Map               func(collector []interface{}) interface{}
}
