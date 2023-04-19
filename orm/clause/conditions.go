package clause

import (
	"fmt"
	"strings"

	"dm.net/datamine/orm/util"
)

// Comparator 条件表达式中的比较符
type Comparator[T interface{}] func(clmn string, val ...T) string

// Equals a=b
func Equals[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s = %s", clmn, util.StringValue(val[0]))
}

// NotEquals a<>b
func NotEquals[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s <> %s", clmn, util.StringValue(val[0]))
}

// GreaterThan a>b
func GreaterThan[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s > %s", clmn, util.StringValue(val[0]))
}

// GreaterEquals a>=b
func GreaterEquals[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s >= %s", clmn, util.StringValue(val[0]))
}

// LessThan a<b
func LessThan[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s < %s", clmn, util.StringValue(val[0]))
}

// LessEquals a<=b
func LessEquals[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s <= %s", clmn, util.StringValue(val[0]))
}

// Like a like 'b'
func Like[T interface{}](clmn string, val ...T) string {
	value := strings.ReplaceAll(util.StringValue(val[0]), "%", "%%")
	return fmt.Sprintf("%s LIKE %s", clmn, value)
}

// NotLike a not like 'b'
func NotLike[T interface{}](clmn string, val ...T) string {
	value := strings.ReplaceAll(util.StringValue(val[0]), "%", "%%")
	return fmt.Sprintf("%s NOT LIKE %s", clmn, value)
}

// Between a between b and c
func Between[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s BETWEEN %s AND %s", clmn, util.StringValue(val[0]), util.StringValue(val[2]))
}

// NotBetween a not between b and c
func NotBetween[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s NOT BETWEEN %s AND %s", clmn, util.StringValue(val[0]), util.StringValue(val[2]))
}

// Is a is b
func Is[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s IS %s", clmn, util.StringValue(val[0]))
}

// IsNot a is not b
func IsNot[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s IS NOT %s", clmn, util.StringValue(val[0]))
}

// In a in (b,c)
func In[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s IN (%s)", clmn, util.StringValues(val))
}

// NotIn a not in (b,c)
func NotIn[T interface{}](clmn string, val ...T) string {
	return fmt.Sprintf("%s NOT IN (%s)", clmn, util.StringValues(val))
}

//ConditionExp 条件
type ConditionExp[VT interface{}] struct {
	Column     string
	Comparator Comparator[VT]
	Value      []VT

	conditions []*ConditionExp[VT] //TODO
	relation   string
}

// Express 条件的字符串化
func (exp *ConditionExp[VT]) Express() string {
	if exp.relation == "" {
		return exp.Comparator(exp.Column, exp.Value...)
	}

	exps := make([]string, len(exp.conditions))
	for i, cond := range exp.conditions {
		exps[i] = cond.Express()
	}
	return strings.Join(exps, exp.relation)
}

// ConditionsJoins 抽象多条件的联合
type ConditionsJoins[T interface{}] func(multicondition ...*ConditionExp[T]) *ConditionExp[T]

// JoinsAndConditions 抽象多条件的 与 联合
func JoinsAndConditions[T interface{}](multicondition ...*ConditionExp[T]) *ConditionExp[T] {
	return &ConditionExp[T]{
		Column:     "",
		Comparator: nil,
		Value:      nil,

		conditions: multicondition,
		relation:   "AND",
	}
}

// JoinsOrConditions 抽象多条件的 或 联合
func JoinsOrConditions(multicondition ...*ConditionExp[interface{}]) *ConditionExp[interface{}] {
	return &ConditionExp[interface{}]{
		Column:     "",
		Comparator: nil,
		Value:      nil,

		conditions: multicondition,
		relation:   "OR",
	}
}

//NewConditionExp 一个条件
func NewConditionExp[T interface{}](column string,
	comparator Comparator[T], value ...T) *ConditionExp[T] {
	return &ConditionExp[T]{Column: column, Comparator: comparator, Value: value}
}
