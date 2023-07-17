package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopBatchTaskRequest Request Object
type StopBatchTaskRequest struct {

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。您可以在IoTDA管理控制台界面，选择左侧导航栏“总览”页签查看当前实例的ID。
	InstanceId *string `json:"Instance-Id,omitempty"`

	// **参数说明**：批量任务ID，创建批量任务时由物联网平台分配获得。 **取值范围**：长度不超过24，只允许小写字母a到f、数字的组合。
	TaskId string `json:"task_id"`

	Body *BatchTargets `json:"body,omitempty"`
}

func (o StopBatchTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopBatchTaskRequest struct{}"
	}

	return strings.Join([]string{"StopBatchTaskRequest", string(data)}, " ")
}
