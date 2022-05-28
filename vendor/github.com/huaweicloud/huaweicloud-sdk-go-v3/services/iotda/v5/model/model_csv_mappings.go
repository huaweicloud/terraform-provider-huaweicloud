package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// csv文件格式转换列表
type CsvMappings struct {

	// **参数说明**：OBS文件中的列名
	ColumnName string `json:"column_name"`

	// **参数说明**：流转数据的属性名
	JsonKey string `json:"json_key"`
}

func (o CsvMappings) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CsvMappings struct{}"
	}

	return strings.Join([]string{"CsvMappings", string(data)}, " ")
}
