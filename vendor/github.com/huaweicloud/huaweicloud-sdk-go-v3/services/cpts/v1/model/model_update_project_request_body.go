package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateProjectRequestBody UpdateProjectRequestBody
type UpdateProjectRequestBody struct {

	// 工程id
	Id int32 `json:"id"`

	// 工程名称
	Name string `json:"name"`

	// 工程描述
	Description *string `json:"description,omitempty"`

	// 导入工程时，缺失的存在于变量文件中的变量
	VariablesNoFile *[]string `json:"variables_no_file,omitempty"`

	// 来源（0-PerfTest；2-CloudTest）
	Source *int32 `json:"source,omitempty"`

	// 扩展参数
	ExternalParams *interface{} `json:"external_params,omitempty"`
}

func (o UpdateProjectRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProjectRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateProjectRequestBody", string(data)}, " ")
}
