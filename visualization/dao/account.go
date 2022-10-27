package dao

import (
	"fmt"
	"log"

	"dm.net/datamine/orm"
	"dm.net/datamine/orm/clause"
	"dm.net/datamine/visualization/model"
	"dm.net/datamine/visualization/resource"
)

func LoadAccounts() []*model.Account {
	conn := resource.GetDBConnection()
	if conn == nil {
		return nil
	}

	var funcs orm.RecordToStructMapping
	funcs.FieldsMapped = func() []string {
		return []string{"*"}
	}
	funcs.DataSourceMapped = func() string { return "t_account" }
	funcs.WhereClause = clause.NewWhereClause(nil, clause.NewConditionExp("phone_number", clause.Equals, "159816559076"))
	funcs.GroupClause = nil
	funcs.OrderClause = nil
	funcs.OffsetLimitClause = nil
	funcs.Collector = func() []interface{} { return nil }
	funcs.Map = func(collector []interface{}) interface{} { return nil }

	result := conn.QueryMultirecord(funcs)
	if result.HasExecuteError() != nil {
		log.Fatal(result.HasExecuteError())
		return nil
	}

	var accounts []*model.Account
	err := result.MapRecords2Struct(&accounts)
	if err != nil {
		fmt.Printf("db error %s", err.Error())
	}
	for _, acc := range accounts {
		fmt.Printf("db object %v", acc)
	}
	return nil
}
