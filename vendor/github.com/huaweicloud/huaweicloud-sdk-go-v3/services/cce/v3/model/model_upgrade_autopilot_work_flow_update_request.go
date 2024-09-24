package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeAutopilotWorkFlowUpdateRequest Request Object
type UpgradeAutopilotWorkFlowUpdateRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 集群升级任务引导流程ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	UpgradeWorkflowId string `json:"upgrade_workflow_id"`

	Body *UpgradeWorkFlowUpdateRequestBody `json:"body,omitempty"`
}

func (o UpgradeAutopilotWorkFlowUpdateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeAutopilotWorkFlowUpdateRequest struct{}"
	}

	return strings.Join([]string{"UpgradeAutopilotWorkFlowUpdateRequest", string(data)}, " ")
}
