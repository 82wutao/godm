package dao

import (
	"fmt"

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
		return []string{"account_id", "name", "avatar", "email", "phone_area", "phone_number", "password"}
	}
	funcs.DataSourceMapped = func() string { return "t_account" }
	funcs.WhereClause = clause.NewWhereClause(nil, clause.NewConditionExp("phone_number", clause.Equals, "159816559076"))
	funcs.GroupClause = nil
	funcs.OrderClause = nil
	funcs.OffsetLimitClause = nil
	funcs.Collector = func() []interface{} { return nil }
	funcs.Map = func(collector []interface{}) interface{} { return nil }

	result := conn.QueryMultirecord(funcs)
	if result.HasError() != nil {
		return nil
	}

	// rcverSupple func() []interface{},
	// structSupple func() interface{},
	// mapFunc func(rcver []interface{}, dest interface{})

	accounts, err := result.MapRecords2Struct(func() []interface{} {
		return []interface{}{new(int64), new(*string), new(*string), new(*string), new(*string), new(*string), new(string)}
	}, func() interface{} {
		return new(model.Account)
	}, func(rcver []interface{}, dest interface{}) {
		acc := dest.(*model.Account)
		acc.AccountID = *rcver[0].(*int64)
		acc.Name = *rcver[1].(**string)
		acc.Avatar = *rcver[2].(**string)
		acc.Email = *rcver[3].(**string)
		acc.PhoneArea = *rcver[4].(**string)
		acc.PhoneNumber = *rcver[5].(**string)
		acc.Password = *rcver[6].(*string)
	})
	if err != nil {
		fmt.Printf("db error %s", err.Error())
	}
	for _, acc := range accounts {
		fmt.Printf("db object %v", acc)
	}
	return nil
}
