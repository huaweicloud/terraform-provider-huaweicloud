package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotAddonInstanceRequest Request Object
type ShowAutopilotAddonInstanceRequest struct {

	// 插件实例id
	Id string `json:"id"`

	// 集群 ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)
	ClusterId *string `json:"cluster_id,omitempty"`
}

func (o ShowAutopilotAddonInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotAddonInstanceRequest struct{}"
	}

	return strings.Join([]string{"ShowAutopilotAddonInstanceRequest", string(data)}, " ")
}
