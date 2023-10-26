package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EsListenerResponse struct {

	// 监听器的监听协议。
	Protocol *string `json:"protocol,omitempty"`

	// 监听器ID。
	Id *string `json:"id,omitempty"`

	// 监听器的名称。
	Name *string `json:"name,omitempty"`

	// 监听器的前端监听端口。
	ProtocolPort *string `json:"protocol_port,omitempty"`

	Ipgroup *EsIpgroupResource `json:"ipgroup,omitempty"`
}

func (o EsListenerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EsListenerResponse struct{}"
	}

	return strings.Join([]string{"EsListenerResponse", string(data)}, " ")
}
