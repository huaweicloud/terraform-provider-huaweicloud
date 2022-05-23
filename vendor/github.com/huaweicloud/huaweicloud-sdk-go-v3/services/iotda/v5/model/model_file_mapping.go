package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OBS文件中的列和流转数据的对应关系。
type FileMapping struct {

	// **参数说明**：csv文件格式转换列表。当file_type为csv时，必填。
	CsvMappings *[]CsvMappings `json:"csv_mappings,omitempty"`
}

func (o FileMapping) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileMapping struct{}"
	}

	return strings.Join([]string{"FileMapping", string(data)}, " ")
}
