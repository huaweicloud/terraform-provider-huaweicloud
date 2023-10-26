package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EsHealthmonitorsResource struct {

	// 后端服务器ID。
	Id *string `json:"id,omitempty"`

	// 后端服务器的名称。
	Name *string `json:"name,omitempty"`

	// 后端服务的前端监听端口。
	ProtocolPort *string `json:"protocol_port,omitempty"`

	Ipgroup *EsHealthIpgroupResource `json:"ipgroup,omitempty"`
}

func (o EsHealthmonitorsResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EsHealthmonitorsResource struct{}"
	}

	return strings.Join([]string{"EsHealthmonitorsResource", string(data)}, " ")
}
