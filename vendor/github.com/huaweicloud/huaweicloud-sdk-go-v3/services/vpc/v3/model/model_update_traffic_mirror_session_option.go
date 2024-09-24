package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorSessionOption
type UpdateTrafficMirrorSessionOption struct {

	// 功能说明：流量镜像会话名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name *string `json:"name,omitempty"`

	// 功能说明：流量镜像会话的描述信息 取值范围：0-255个字符，不能包含“<”和“>”
	Description *string `json:"description,omitempty"`

	// 功能说明：流量镜像筛选条件ID
	TrafficMirrorFilterId *string `json:"traffic_mirror_filter_id,omitempty"`

	// 功能说明：镜像目标ID
	TrafficMirrorTargetId *string `json:"traffic_mirror_target_id,omitempty"`

	// 功能说明：镜像目的类型 取值范围：     eni：弹性网卡     elb：私网弹性负载均衡
	TrafficMirrorTargetType *string `json:"traffic_mirror_target_type,omitempty"`

	// 功能说明：指定VNI，用于在镜像目的区分不同会话的镜像流量 取值范围：0~16777215
	VirtualNetworkId *int32 `json:"virtual_network_id,omitempty"`

	// 功能说明：最大传输单元MTU 取值范围：1~1460
	PacketLength *int32 `json:"packet_length,omitempty"`

	// 功能说明：会话优先级 取值范围：1~32766
	Priority *int32 `json:"priority,omitempty"`

	// 功能说明：是否开启会话 取值范围：true、false
	Enabled *string `json:"enabled,omitempty"`
}

func (o UpdateTrafficMirrorSessionOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorSessionOption struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorSessionOption", string(data)}, " ")
}
