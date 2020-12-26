package syntaxutil

// TernaryOperate 常见的三目运算
func TernaryOperate(t bool, trueValue interface{}, falseValue func() interface{}) interface{} {
	if t {
		return trueValue
	}
	return falseValue()
}
