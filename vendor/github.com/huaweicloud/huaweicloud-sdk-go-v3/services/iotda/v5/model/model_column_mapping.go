package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ColumnMapping 数据库的列和流转数据的对应关系。
type ColumnMapping struct {

	// **参数说明**：数据库的列名
	ColumnName string `json:"column_name"`

	// **参数说明**：流转数据的属性名
	JsonKey string `json:"json_key"`
}

func (o ColumnMapping) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ColumnMapping struct{}"
	}

	return strings.Join([]string{"ColumnMapping", string(data)}, " ")
}
