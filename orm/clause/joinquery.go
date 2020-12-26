package clause

import (
	"fmt"
	"strings"

	"dm.net/datamine/syntaxutil"
)

// Joinpointer 联合查询的联合点
type Joinpointer struct {
	ATblKey   string
	BTblKey   string
	ATblAlias string
	BTblAlias string
}

// JoinQuery 多表联合查询
type JoinQuery interface {
	DataSourceMapped() TargetDataSource
	InnerJoin(otherTbl, otherTblAlias string, pointers ...Joinpointer) JoinQuery
	LeftJoin(otherTbl, otherTblAlias string, pointers ...Joinpointer) JoinQuery
	RightJoin(otherTbl, otherTblAlias string, pointers ...Joinpointer) JoinQuery
	FullJoin(otherTbl, otherTblAlias string, pointers ...Joinpointer) JoinQuery
}

// NewMultitableJoinQuery 新建一个多表联合查询
func NewMultitableJoinQuery(mainTbl, mainTblAlias string) JoinQuery {
	return &simpleJoinQueryImpl{mainTbl: mainTbl, mainAlias: mainTblAlias}
}

type simpleJoinpointerImpl struct {
	mainTblAlias string
	joinWay      string
	tbl          string
	tblAlias     string
	pointers     []Joinpointer
}

func (impl *simpleJoinpointerImpl) Express() string {

	pointers := make([]string, len(impl.pointers))
	for i, p := range impl.pointers {
		aTbl := syntaxutil.TernaryOperate(p.ATblAlias == "", impl.mainTblAlias, func() interface{} { return p.ATblAlias })
		bTbl := syntaxutil.TernaryOperate(p.BTblAlias == "", impl.tblAlias, func() interface{} { return p.BTblAlias })
		pointers[i] = fmt.Sprintf("%s.%s = %s.%s", aTbl, p.ATblKey, bTbl, p.BTblKey)
	}
	pointersSQL := strings.Join(pointers, " AND ")

	return fmt.Sprintf("%s %s %s ON %s",
		impl.joinWay,
		impl.tbl, impl.tblAlias,
		pointersSQL)
}

type simpleJoinQueryImpl struct {
	mainTbl   string
	mainAlias string
	moreTbls  []*simpleJoinpointerImpl
}

func (joinQuery *simpleJoinQueryImpl) DataSourceMapped() TargetDataSource {
	return func() string {
		tbls := make([]string, len(joinQuery.moreTbls)+1)
		tbls[0] = fmt.Sprintf("%s %s", joinQuery.mainTbl, joinQuery.mainAlias)

		for i, t := range joinQuery.moreTbls {
			tbls[i+1] = t.Express()
		}
		return strings.Join(tbls, " ")
	}
}

func (joinQuery *simpleJoinQueryImpl) InnerJoin(otherTbl, otherTblAlias string, pointers ...Joinpointer) JoinQuery {
	pointer := &simpleJoinpointerImpl{
		mainTblAlias: joinQuery.mainAlias,
		joinWay:      "INNER JOIN",
		tbl:          otherTbl,
		tblAlias:     otherTblAlias,
		pointers:     pointers}

	joinQuery.moreTbls = append(joinQuery.moreTbls, pointer)
	return joinQuery
}
func (joinQuery *simpleJoinQueryImpl) LeftJoin(otherTbl, otherTblAlias string, pointers ...Joinpointer) JoinQuery {
	pointer := &simpleJoinpointerImpl{
		mainTblAlias: joinQuery.mainAlias,
		joinWay:      "LEFT OUTER JOIN",
		tbl:          otherTbl,
		tblAlias:     otherTblAlias,
		pointers:     pointers}

	joinQuery.moreTbls = append(joinQuery.moreTbls, pointer)
	return joinQuery
}
func (joinQuery *simpleJoinQueryImpl) RightJoin(otherTbl, otherTblAlias string, pointers ...Joinpointer) JoinQuery {
	pointer := &simpleJoinpointerImpl{
		mainTblAlias: joinQuery.mainAlias,
		joinWay:      "RIGHT OUTER JOIN",
		tbl:          otherTbl,
		tblAlias:     otherTblAlias,
		pointers:     pointers}

	joinQuery.moreTbls = append(joinQuery.moreTbls, pointer)
	return joinQuery
}
func (joinQuery *simpleJoinQueryImpl) FullJoin(otherTbl, otherTblAlias string, pointers ...Joinpointer) JoinQuery {
	pointer := &simpleJoinpointerImpl{
		mainTblAlias: joinQuery.mainAlias,
		joinWay:      "FULL OUTER JOIN",
		tbl:          otherTbl,
		tblAlias:     otherTblAlias,
		pointers:     pointers}

	joinQuery.moreTbls = append(joinQuery.moreTbls, pointer)
	return joinQuery
}
