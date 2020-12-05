package orm

import "errors"

var (
	ERR_CREATE_PARAM_MUST_NOT_BE_EMPTY     = errors.New("创建或修改时填充参数不能为空")
	ERR_CREATE_LEN_FIELDS_MUST_SAME_VALUES = errors.New("创建或修改时填充的参数必须等于填充的字段")
	ERR_REFLECT_DATAKIND_CANOT_BE_REFLECT  = errors.New("填充参数时数据类型不能被反射")
	ERR_QUERY_SELECT_CANOT_BE_BLANK        = errors.New("查询时选中字段不能为空")
	ERR_ORM_STRUCT_MUST_OVERWRITE_STRING   = errors.New("ORM映射时结构没有实现String方法")
	ERR_CANOT_MAP_RECORD_2_MEM             = errors.New("")
	ERR_NOT_VALIDATE_ARGUMENT              = errors.New("")
	ERR_STORE_CONTAINER_IS_EMPTY           = errors.New("")
	ERR_STORE_ELE_NOT_POINTER              = errors.New("")
)
