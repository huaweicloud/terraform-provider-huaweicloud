package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TaskBasicAttribute struct {

	// 分支ID
	BranchId *string `json:"branch_id,omitempty"`

	// 分支名
	BranchName *string `json:"branch_name,omitempty"`

	// 创建人的工号
	CreateBy *string `json:"create_by,omitempty"`

	// 迭代url
	IterationUri *string `json:"iteration_uri,omitempty"`

	// 工程id
	ProjectId *string `json:"project_id,omitempty"`

	// 协议
	Protocols *[]string `json:"protocols,omitempty"`

	// 服务id
	ServiceId *string `json:"service_id,omitempty"`

	// 阶段
	Stage *int32 `json:"stage,omitempty"`

	// 阶段名称
	StageName *string `json:"stage_name,omitempty"`

	// 任务id
	TaskId *string `json:"task_id,omitempty"`

	// 版本uri
	VersionUri *string `json:"version_uri,omitempty"`
}

func (o TaskBasicAttribute) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskBasicAttribute struct{}"
	}

	return strings.Join([]string{"TaskBasicAttribute", string(data)}, " ")
}
