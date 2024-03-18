package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTempRequestBody UpdateTempRequestBody
type UpdateTempRequestBody struct {

	// 事务id
	Id int32 `json:"id"`

	// 工程id
	ProjectId int32 `json:"project_id"`

	// 事务名称
	Name string `json:"name"`

	// 事务类型
	TempType *int32 `json:"temp_type,omitempty"`

	// 描述信息
	Description *string `json:"description,omitempty"`

	// 旧版本逻辑控制器字段，当前已未使用
	ForLoopParams *[]interface{} `json:"for_loop_params,omitempty"`

	// 是否启用预置事务，当前版本已未使用
	EnablePre *bool `json:"enable_pre,omitempty"`

	// 事务脚本信息
	Contents *[]TempContentInfo `json:"contents,omitempty"`
}

func (o UpdateTempRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTempRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateTempRequestBody", string(data)}, " ")
}
