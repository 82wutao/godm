package clause

// TargetFields 被映射的字段集合
type TargetFields func() []string

//TargetDataSource 被映射的数据源
type TargetDataSource func() string
