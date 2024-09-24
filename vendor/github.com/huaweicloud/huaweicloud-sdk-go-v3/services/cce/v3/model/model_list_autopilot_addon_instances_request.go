package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotAddonInstancesRequest Request Object
type ListAutopilotAddonInstancesRequest struct {

	// 含义：想要筛选的插件名称或插件别名  属性：隐藏参数
	AddonTemplateName *string `json:"addon_template_name,omitempty"`

	// 集群 ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)
	ClusterId string `json:"cluster_id"`
}

func (o ListAutopilotAddonInstancesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotAddonInstancesRequest struct{}"
	}

	return strings.Join([]string{"ListAutopilotAddonInstancesRequest", string(data)}, " ")
}
