package clause

import "fmt"

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
