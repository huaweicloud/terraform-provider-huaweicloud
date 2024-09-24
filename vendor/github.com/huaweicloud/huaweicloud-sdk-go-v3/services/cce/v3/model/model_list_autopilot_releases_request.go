package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotReleasesRequest Request Object
type ListAutopilotReleasesRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 模板ID
	ChartId *string `json:"chart_id,omitempty"`

	// 模板对应的命名空间
	Namespace *string `json:"namespace,omitempty"`
}

func (o ListAutopilotReleasesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotReleasesRequest struct{}"
	}

	return strings.Join([]string{"ListAutopilotReleasesRequest", string(data)}, " ")
}
