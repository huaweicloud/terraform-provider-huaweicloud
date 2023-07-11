package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NewTaskInfo struct {

	// 任务名称
	Name string `json:"name"`

	// 是否并行
	Parallel bool `json:"parallel"`

	// 是否支持全链路压测
	EnableFullLink *bool `json:"enable_full_link,omitempty"`

	// 工程id
	ProjectId int32 `json:"project_id"`

	// 任务压测模式，0-时长模式；1-次数模式；2-混合模式；此处是兼容老版本遗留字段，填固定值2
	OperateMode int32 `json:"operate_mode"`

	// 关联用例id列表
	CaseIdList []int32 `json:"case_id_list"`
}

func (o NewTaskInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NewTaskInfo struct{}"
	}

	return strings.Join([]string{"NewTaskInfo", string(data)}, " ")
}
