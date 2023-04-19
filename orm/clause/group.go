package clause

import (
	"fmt"
	"strings"
)

//GroupClause 分组子句
type GroupClause interface {
	GroupSQL() string
}
type simpleGroupClauseImpl[T interface{}] struct {
	columns    []string
	conditions []*ConditionExp[T]
	joins      ConditionsJoins[T]
}

//Express 分组子句的sql表达
func (group *simpleGroupClauseImpl[T]) GroupSQL() string {

	clmns := strings.Join(group.columns, ",")
	groupSQL := fmt.Sprintf("GROUP BY %s", clmns)
	if len(group.conditions) == 0 {
		return groupSQL
	}

	joins := group.joins
	if joins == nil {
		joins = func(multicondition ...*ConditionExp[T]) *ConditionExp[T] {
			return multicondition[0]
		}
	}
	return fmt.Sprintf("%s HAVING %s", groupSQL, joins(group.conditions...).Express())
}

//NewGroupClause group子句
func NewGroupClause[T interface{}](columns []string,
	havings []*ConditionExp[T], joins ConditionsJoins[T]) GroupClause {
	return &simpleGroupClauseImpl[T]{columns: columns, conditions: havings, joins: joins}
}
