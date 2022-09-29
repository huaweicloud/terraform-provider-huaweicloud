package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 节点。
type ShowClusterRespClusterBrokers struct {

	// 节点IP。
	Host *string `json:"host,omitempty"`

	// 端口号。
	Port *int32 `json:"port,omitempty"`

	// 节点ID。
	BrokerId *string `json:"broker_id,omitempty"`

	// 是否为contoller节点。
	IsController *bool `json:"is_controller,omitempty"`

	// 服务端版本。
	Version *string `json:"version,omitempty"`

	// broker注册时间，为unix时间戳格式。
	RegisterTime *int64 `json:"register_time,omitempty"`

	// Kafka实例节点的连通性是否正常。
	IsHealth *bool `json:"is_health,omitempty"`
}

func (o ShowClusterRespClusterBrokers) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterRespClusterBrokers struct{}"
	}

	return strings.Join([]string{"ShowClusterRespClusterBrokers", string(data)}, " ")
}
