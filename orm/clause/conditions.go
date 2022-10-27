package clause

import (
	"fmt"
	"strings"

	"dm.net/datamine/orm/util"
)

// Comparator 条件表达式中的比较符
type Comparator func(clmn string, val ...interface{}) string

var (
	// Equals a=b
	Equals Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s = %s", clmn, util.StringValue(val[0]))
	}
	// NotEquals a<>b
	NotEquals Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s <> %s", clmn, util.StringValue(val[0]))
	}
	// GreaterThan a>b
	GreaterThan Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s > %s", clmn, util.StringValue(val[0]))
	}
	// GreaterEquals a>=b
	GreaterEquals Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s >= %s", clmn, util.StringValue(val[0]))
	}
	// LessThan a<b
	LessThan Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s < %s", clmn, util.StringValue(val[0]))
	}
	// LessEquals a<=b
	LessEquals Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s <= %s", clmn, util.StringValue(val[0]))
	}
	// Like a like 'b'
	Like Comparator = func(clmn string, val ...interface{}) string {
		value := strings.ReplaceAll(util.StringValue(val[0]), "%", "%%")
		return fmt.Sprintf("%s LIKE %s", clmn, value)
	}
	// NotLike a not like 'b'
	NotLike Comparator = func(clmn string, val ...interface{}) string {
		value := strings.ReplaceAll(util.StringValue(val[0]), "%", "%%")
		return fmt.Sprintf("%s NOT LIKE %s", clmn, value)
	}
	// Between a between b and c
	Between Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s BETWEEN %s AND %s", clmn, util.StringValue(val[0]), util.StringValue(val[2]))
	}
	// NotBetween a not between b and c
	NotBetween Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s NOT BETWEEN %s AND %s", clmn, util.StringValue(val[0]), util.StringValue(val[2]))
	}
	// Is a is b
	Is Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s IS %s", clmn, util.StringValue(val[0]))
	}
	// IsNot a is not b
	IsNot Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s IS NOT %s", clmn, util.StringValue(val[0]))
	}
	// In a in (b,c)
	In Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s IN (%s)", clmn, util.StringValues(val))
	}
	// NotIn a not in (b,c)
	NotIn Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s NOT IN (%s)", clmn, util.StringValues(val))
	}
)

//ConditionExp 条件
type ConditionExp struct {
	Column     string
	Comparator Comparator
	Value      []interface{}

	conditions []*ConditionExp
	relation   string
}

// Express 条件的字符串化
func (exp *ConditionExp) Express() string {
	if exp.relation == "" {
		return exp.Comparator(exp.Column, exp.Value)
	}

	exps := make([]string, len(exp.conditions))
	for i, cond := range exp.conditions {
		exps[i] = cond.Express()
	}
	return strings.Join(exps, exp.relation)
}

// ConditionsJoins 抽象多条件的联合
type ConditionsJoins func(multicondition ...*ConditionExp) *ConditionExp

// JoinsAndConditions 抽象多条件的 与 联合
func JoinsAndConditions(multicondition ...*ConditionExp) *ConditionExp {
	return &ConditionExp{
		Column:     "",
		Comparator: nil,
		Value:      nil,

		conditions: multicondition,
		relation:   "AND",
	}
}

// JoinsOrConditions 抽象多条件的 或 联合
func JoinsOrConditions(multicondition ...*ConditionExp) *ConditionExp {
	return &ConditionExp{
		Column:     "",
		Comparator: nil,
		Value:      nil,

		conditions: multicondition,
		relation:   "OR",
	}
}

//NewConditionExp 一个条件
func NewConditionExp(column string, comparator Comparator, value ...interface{}) *ConditionExp {
	return &ConditionExp{Column: column, Comparator: comparator, Value: value}
}
