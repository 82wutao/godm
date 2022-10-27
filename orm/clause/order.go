package clause

import (
	"fmt"
	"strings"

	"dm.net/datamine/orm/util"
)

//OrderClause 排序
type OrderClause interface {
	AddASCFactor(column string)
	AddDESCFactor(column string)
	OrderSQL() string
}

type orderFactor struct {
	column string
	order  string
}

func (f *orderFactor) String() string {
	return fmt.Sprintf("%s %s", f.column, f.order)
}

type simpleOrderClauseImpl struct {
	factors []orderFactor
}

func (oc *simpleOrderClauseImpl) OrderSQL() string {
	if len(oc.factors) == 0 {
		return ""
	}
	return fmt.Sprintf("ORDER BY %s", util.StringValue(oc.factors))
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
