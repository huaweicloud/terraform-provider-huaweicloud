package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchUpdateTaskStatusRequestBody struct {

	// 类型（0-旧版本任务；1-新版本任务）
	Type *int32 `json:"type,omitempty"`

	// 所属工程id
	ProjectId int32 `json:"project_id"`

	// 任务id列表
	TaskIdList []int32 `json:"task_id_list"`

	// cluster_id
	ClusterId int32 `json:"cluster_id"`

	// cluster_type
	ClusterType string `json:"cluster_type"`

	// without_package
	WithoutPackage *int32 `json:"without_package,omitempty"`

	NetworkInfo *NetworkInfo `json:"network_info,omitempty"`

	// status
	Status int32 `json:"status"`

	// 企业项目id
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o BatchUpdateTaskStatusRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchUpdateTaskStatusRequestBody struct{}"
	}

	return strings.Join([]string{"BatchUpdateTaskStatusRequestBody", string(data)}, " ")
}
