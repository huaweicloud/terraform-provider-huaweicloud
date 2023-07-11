package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateNewTaskRequestBody 新版(任务用例分离版本)更新任务请求体
type UpdateNewTaskRequestBody struct {

	// 名称
	Name string `json:"name"`

	// 并行状态（表示任务执行时用例是否并行执行；true：并行执行，false：串行执行）
	Parallel bool `json:"parallel"`

	// 工程id
	ProjectId int32 `json:"project_id"`

	// 任务模式（兼容旧版接口保留字段，0：时长模式，1：次数模式，2：混合模式；此处请传固定值：2）
	OperateMode int32 `json:"operate_mode"`

	// 关联的用例id集合
	CaseIdList []int32 `json:"case_id_list"`
}

func (o UpdateNewTaskRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateNewTaskRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateNewTaskRequestBody", string(data)}, " ")
}
