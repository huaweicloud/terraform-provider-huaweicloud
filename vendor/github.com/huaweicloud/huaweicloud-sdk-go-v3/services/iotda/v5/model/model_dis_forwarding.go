package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DIS服务配置信息
type DisForwarding struct {

	// **参数说明**：DIS服务对应的region区域
	RegionName string `json:"region_name"`

	// **参数说明**：DIS服务对应的projectId信息
	ProjectId string `json:"project_id"`

	// **参数说明**：DIS服务对应的通道名称，stream_id和stream_name两个参数必须携带一个，优先使用stream_id
	StreamName *string `json:"stream_name,omitempty"`

	// **参数说明**：DIS服务对应的通道ID，stream_id和stream_name两个参数必须携带一个，优先使用stream_id
	StreamId *string `json:"stream_id,omitempty"`
}

func (o DisForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DisForwarding struct{}"
	}

	return strings.Join([]string{"DisForwarding", string(data)}, " ")
}
