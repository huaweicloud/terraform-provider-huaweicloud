package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTempRequestBody CreateTempRequestBody
type CreateTempRequestBody struct {

	// 所属工程id
	ProjectId int32 `json:"project_id"`

	// 事务类型
	TempType int32 `json:"temp_type"`

	// 事务名称
	Name string `json:"name"`

	// 描述信息
	Description *string `json:"description,omitempty"`

	// 事务脚本信息
	Contents *[]interface{} `json:"contents,omitempty"`
}

func (o CreateTempRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTempRequestBody struct{}"
	}

	return strings.Join([]string{"CreateTempRequestBody", string(data)}, " ")
}
