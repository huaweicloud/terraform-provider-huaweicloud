package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TempInfo struct {

	// 事务id
	Id *int32 `json:"id,omitempty"`

	// 工程id
	ProjectId *int32 `json:"project_id,omitempty"`

	// 事务名称
	Name *string `json:"name,omitempty"`

	// 事务描述
	Description *string `json:"description,omitempty"`

	// 变量
	Variables *string `json:"variables,omitempty"`

	// 事务脚本信息
	Contents *[]interface{} `json:"contents,omitempty"`

	// 事务类型（弃用）
	TempType *int32 `json:"temp_type,omitempty"`

	// 旧版本逻辑控制器字段，当前已未使用
	ForLoopParams *[]interface{} `json:"for_loop_params,omitempty"`

	LogicController *LogicController `json:"logic_controller,omitempty"`

	// 是否启用预置事务，当前版本已未使用
	EnablePre *bool `json:"enable_pre,omitempty"`
}

func (o TempInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TempInfo struct{}"
	}

	return strings.Join([]string{"TempInfo", string(data)}, " ")
}
