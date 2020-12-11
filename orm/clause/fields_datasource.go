package clause

// FieldsMapped 被映射的字段集合
type FieldsMapped interface {
	Fields() []string
}

//DataSourceMapped 被映射的数据源
type DataSourceMapped interface {
	DataSource() string
}
