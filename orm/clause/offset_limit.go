package clause

import (
	"fmt"
)

//OffsetLimitClause 偏移量子句
type OffsetLimitClause interface {
	OffsetLimitSQL() string
}
type simpleOffsetLimitClauseImpl struct {
	am     string
	offset int
	limit  int
}

//Express 分组子句的sql表达
func (offsetLimit *simpleOffsetLimitClauseImpl) OffsetLimitSQL() string {
	switch offsetLimit.am {
	case "ORACLE":
	case "MYSQL":
		return fmt.Sprintf("LIMIT %d , %d", offsetLimit.offset, offsetLimit.limit)
	case "SQLSERVER":
		return fmt.Sprintf("OFFSET %d ROW FETCH NEXT %d ROWS ONLY", offsetLimit.offset, offsetLimit.limit)
	case "POSTGRESQL":
		return fmt.Sprintf("OFFSET %d LIMIT %d", offsetLimit.offset, offsetLimit.limit)
	}
	return fmt.Sprintf("OFFSET %d LIMIT %d", offsetLimit.offset, offsetLimit.limit)
}

//NewOracleOffsetLimitClause 偏移量子句
func NewOracleOffsetLimitClause(offset, limit int) OffsetLimitClause {
	return &simpleOffsetLimitClauseImpl{offset: offset, limit: limit}
}

//NewMySQLOffsetLimitClause 偏移量子句
func NewMySQLOffsetLimitClause(offset, limit int) OffsetLimitClause {
	return &simpleOffsetLimitClauseImpl{offset: offset, limit: limit}
}

//NewSQLServerOffsetLimitClause 偏移量子句
func NewSQLServerOffsetLimitClause(offset, limit int) OffsetLimitClause {
	return &simpleOffsetLimitClauseImpl{offset: offset, limit: limit}
}

//NewPostgreSQLOffsetLimitClause 偏移量子句
func NewPostgreSQLOffsetLimitClause(offset, limit int) OffsetLimitClause {
	return &simpleOffsetLimitClauseImpl{offset: offset, limit: limit}
}
