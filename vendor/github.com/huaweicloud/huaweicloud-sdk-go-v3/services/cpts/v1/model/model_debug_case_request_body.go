package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DebugCaseRequestBody DebugCaseRequestBody
type DebugCaseRequestBody struct {

	// 状态（9：启动调试）
	Status int32 `json:"status"`

	// 资源组id
	ClusterId int32 `json:"cluster_id"`

	// 资源组类型（共享资源组：shared-cluster-internet；私有资源组：private-cluster）
	ClusterType string `json:"cluster_type"`

	// 套餐包VUM不足的情况下用户选择是不是要走按需计费模式（当前版本固定值：0）
	WithoutPackage int32 `json:"without_package"`
}

func (o DebugCaseRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseRequestBody struct{}"
	}

	return strings.Join([]string{"DebugCaseRequestBody", string(data)}, " ")
}
