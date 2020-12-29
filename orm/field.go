package orm

import (
	"fmt"
	"reflect"
	"strings"

	"dm.net/datamine/progmode/stream"
)

type FieldMapping struct {
	StructField string
	FieldType   reflect.Type
	FieldCache  reflect.Value
	ColumnName  string
	ColumnType  string
}

func parseColumnMapping(ormStr string) *FieldMapping {
	//gorm:"column:phone_number;type:varchar(32);default:null"
	//TODO complete feature

	var mapping FieldMapping
	kvArr := strings.Split(ormStr, ";")
	for _, kv := range kvArr {
		kv := strings.Split(kv, ":")
		switch kv[0] {
		case "column":
			mapping.ColumnName = kv[1]
		case "type":
			mapping.ColumnType = kv[1]
		}
	}
	return &mapping
}

func findFields(structType reflect.Type, targetFields ...string) map[string]*FieldMapping {

	selected := stream.NewStreamFromSlice(targetFields).Group(func() interface{} {
		return make(map[string]int)
	}, func(src interface{}) (interface{}, interface{}) {
		return src, 1
	}).(map[string]int)

	cout := structType.NumField()

	var ret = make(map[string]*FieldMapping)
	for i := 0; i < cout; i++ {
		field := structType.Field(i)
		_, isAll := selected["*"]
		_, nameMatch := selected[field.Name]
		if !isAll && !nameMatch {
			continue
		}

		ormStr, ok := field.Tag.Lookup("gorm")
		if !ok {
			ormStr = fmt.Sprintf("column:%s", field.Name)
		}
		if ormStr == "-" {
			continue
		}
		mapping := parseColumnMapping(ormStr)
		mapping.StructField = field.Name
		mapping.FieldType = field.Type
		mapping.FieldCache = reflect.New(field.Type).Elem()

		//TODO 比如聚合
		ret[mapping.ColumnName] = mapping
	}
	return ret
}
