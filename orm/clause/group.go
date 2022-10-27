package clause

import (
	"fmt"
	"strings"
)

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
