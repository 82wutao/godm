package orm

import (
	"fmt"
	"strings"
)

type SelectClause interface {
	SelectSQL() string
}
type FromClause interface {
	FromSQL() string
}

//#################

//WhereClause 条件子句
type WhereClause interface {
	WhereSQL() string
}
type simpleWhereClauseImpl struct {
	joins      ConditionsJoins
	conditions []*ConditionExp
}

//Express 条件子句 的sql表达
func (where *simpleWhereClauseImpl) WhereSQL() string {
	if len(where.conditions) == 0 {
		return ""
	}

	joins := where.joins
	if joins == nil {
		joins = func(multicondition ...*ConditionExp) *ConditionExp {
			return multicondition[0]
		}
	}
	return fmt.Sprintf("WHERE %s", joins(where.conditions...).Express())
}

//NewWhereClause where子句
func NewWhereClause(joins ConditionsJoins, conditions ...*ConditionExp) WhereClause {
	return &simpleWhereClauseImpl{joins: joins, conditions: conditions}
}

//#################

//GroupClause 分组子句
type GroupClause interface {
	GroupSQL() string
}
type simpleGroupClauseImpl struct {
	columns    []string
	conditions []*ConditionExp
	joins      ConditionsJoins
}

//Express 分组子句的sql表达
func (group *simpleGroupClauseImpl) GroupSQL() string {

	clmns := strings.Join(group.columns, ",")
	groupSQL := fmt.Sprintf("GROUP BY %s", clmns)
	if len(group.conditions) == 0 {
		return groupSQL
	}

	joins := group.joins
	if joins == nil {
		joins = func(multicondition ...*ConditionExp) *ConditionExp {
			return multicondition[0]
		}
	}
	return fmt.Sprintf("%s HAVING %s", groupSQL, joins(group.conditions...).Express())
}

//NewGroupClause group子句
func NewGroupClause(columns []string, havings []*ConditionExp, joins ConditionsJoins) GroupClause {
	return &simpleGroupClauseImpl{columns: columns, conditions: havings, joins: joins}
}

//#################

//OrderClause 排序
type OrderClause interface {
	AddASCFactor(column string)
	AddDESCFactor(column string)
	OrderSQL() string
}

type simpleOrderClauseImpl struct {
	factors []orderFactor
}

func (oc *simpleOrderClauseImpl) OrderSQL() string {
	if len(oc.factors) == 0 {
		return ""
	}
	return fmt.Sprintf("ORDER BY %s", stringValue(oc.factors))
}
func (oc *simpleOrderClauseImpl) AddASCFactor(column string) {
	factor := orderFactor{
		column: column, order: "ASC",
	}
	oc.factors = append(oc.factors, factor)
}
func (oc *simpleOrderClauseImpl) AddDESCFactor(column string) {
	factor := orderFactor{
		column: column, order: "DESC",
	}
	oc.factors = append(oc.factors, factor)
}

// 偏移，限制
// 分组
// 排序
// union
// 值和类型反射
// 事务 aop

// #where
// #group
// #order
// offset ,limit

// subquery

// Comparator 条件表达式中的比较符
type Comparator func(clmn string, val ...interface{}) string

var (
	// Equals a=b
	Equals Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s = %s", clmn, stringValue(val[0]))
	}
	// NotEquals a<>b
	NotEquals Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s <> %s", clmn, stringValue(val[0]))
	}
	// GreaterThan a>b
	GreaterThan Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s > %s", clmn, stringValue(val[0]))
	}
	// GreaterEquals a>=b
	GreaterEquals Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s >= %s", clmn, stringValue(val[0]))
	}
	// LessThan a<b
	LessThan Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s < %s", clmn, stringValue(val[0]))
	}
	// LessEquals a<=b
	LessEquals Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s <= %s", clmn, stringValue(val[0]))
	}
	// Like a like 'b'
	Like Comparator = func(clmn string, val ...interface{}) string {
		value := strings.ReplaceAll(stringValue(val[0]), "%", "%%")
		return fmt.Sprintf("%s LIKE %s", clmn, value)
	}
	// NotLike a not like 'b'
	NotLike Comparator = func(clmn string, val ...interface{}) string {
		value := strings.ReplaceAll(stringValue(val[0]), "%", "%%")
		return fmt.Sprintf("%s NOT LIKE %s", clmn, value)
	}
	// Between a between b and c
	Between Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s BETWEEN %s AND %s", clmn, stringValue(val[0]), stringValue(val[2]))
	}
	// NotBetween a not between b and c
	NotBetween Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s NOT BETWEEN %s AND %s", clmn, stringValue(val[0]), stringValue(val[2]))
	}
	// Is a is b
	Is Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s IS %s", clmn, stringValue(val[0]))
	}
	// IsNot a is not b
	IsNot Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s IS NOT %s", clmn, stringValue(val[0]))
	}
	// In a in (b,c)
	In Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s IN (%s)", clmn, stringValues(val))
	}
	// NotIn a not in (b,c)
	NotIn Comparator = func(clmn string, val ...interface{}) string {
		return fmt.Sprintf("%s NOT IN (%s)", clmn, stringValues(val))
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

//NewOrderClause order子句
func NewOrderClause(column, order string) OrderClause {
	if order == "" {
		order = "ASC"
	}
	if strings.ToUpper(order) == "ASC" {
		order = "ASC"
	} else {
		order = "DESC"
	}
	factor := orderFactor{column: column, order: order}

	clause := &simpleOrderClauseImpl{}
	clause.factors = make([]orderFactor, 0)
	clause.factors = append(clause.factors, factor)
	return clause
}

type orderFactor struct {
	column string
	order  string
}

func (f *orderFactor) String() string {
	return fmt.Sprintf("%s %s", f.column, f.order)
}

type JoinQuery interface {
	FromClause
	InnerJoin(tbl, alias string) Joinpointer
	LeftJoin(tbl, alias string) Joinpointer
	RightJoin(tbl, alias string) Joinpointer
	FullJoin(tbl, alias string) Joinpointer
}
type Joinpointer interface {
	On(mainKey, thisKey string) JoinQuery
	Express() string
}

func NewMultitableQuery(mainTbl, mainAlias string) JoinQuery {
	return &simpleJoinQueryImpl{mainTbl: mainTbl, mainAlias: mainAlias}
}

type simpleJoinQueryImpl struct {
	mainTbl   string
	mainAlias string
	moreTbls  []Joinpointer
}
type simpleJoinpointerImpl struct {
	joinWay  string
	tbl      string
	tblAlias string
	mainKey  string
	thisKey  string

	mainTbl *simpleJoinQueryImpl
}

func (pointer *simpleJoinpointerImpl) On(mainKey, thisKey string) JoinQuery {
	pointer.mainKey = mainKey
	pointer.thisKey = thisKey
	return pointer.mainTbl
}
func (pointer *simpleJoinpointerImpl) Express() string {
	return fmt.Sprintf("%s %s %s ON %s.%s = %s.%s",
		pointer.joinWay,
		pointer.tbl, pointer.tblAlias,
		pointer.mainTbl.mainAlias, pointer.mainKey, pointer.tblAlias, pointer.thisKey)
}

func (joinQuery *simpleJoinQueryImpl) FromSQL() string {
	tbls := make([]string, len(joinQuery.moreTbls)+1)
	tbls[0] = fmt.Sprintf("%s %s", joinQuery.mainTbl, joinQuery.mainAlias)

	for i, t := range joinQuery.moreTbls {
		tbls[i+1] = t.Express()
	}
	return strings.Join(tbls, " ")
}
func (joinQuery *simpleJoinQueryImpl) InnerJoin(tbl, alias string) Joinpointer {
	pointer := &simpleJoinpointerImpl{joinWay: "INNER JOIN", tbl: tbl, tblAlias: alias}
	joinQuery.moreTbls = append(joinQuery.moreTbls, pointer)
	return pointer
}
func (joinQuery *simpleJoinQueryImpl) LeftJoin(tbl, alias string) Joinpointer {
	pointer := &simpleJoinpointerImpl{joinWay: "LEFT OUTER JOIN", tbl: tbl, tblAlias: alias}
	joinQuery.moreTbls = append(joinQuery.moreTbls, pointer)
	return pointer
}
func (joinQuery *simpleJoinQueryImpl) RightJoin(tbl, alias string) Joinpointer {
	pointer := &simpleJoinpointerImpl{joinWay: "RIGHT OUTER JOIN", tbl: tbl, tblAlias: alias}
	joinQuery.moreTbls = append(joinQuery.moreTbls, pointer)
	return pointer
}
func (joinQuery *simpleJoinQueryImpl) FullJoin(tbl, alias string) Joinpointer {
	pointer := &simpleJoinpointerImpl{joinWay: "FULL OUTER JOIN", tbl: tbl, tblAlias: alias}
	joinQuery.moreTbls = append(joinQuery.moreTbls, pointer)
	return pointer
}

// select *
// from t
// left join t2 on ....
// where
// group by
// order by
// offset limit
