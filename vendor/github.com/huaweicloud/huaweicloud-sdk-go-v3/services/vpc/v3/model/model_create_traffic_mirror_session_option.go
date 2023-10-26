package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorSessionOption
type CreateTrafficMirrorSessionOption struct {

	// 功能说明：流量镜像会话名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：流量镜像会话的描述信息 取值范围：0-255个字符，不能包含“<”和“>”
	Description *string `json:"description,omitempty"`

	// 功能说明：流量镜像筛选条件ID
	TrafficMirrorFilterId string `json:"traffic_mirror_filter_id"`

	// 功能说明：镜像源ID列表，支持弹性网卡作为镜像源。 约束：一个镜像会话默认最大支持10个镜像源。
	TrafficMirrorSources []string `json:"traffic_mirror_sources"`

	// 功能说明：镜像目的ID
	TrafficMirrorTargetId string `json:"traffic_mirror_target_id"`

	// 功能说明：镜像目的类型 取值范围：     eni：弹性网卡     elb：私网弹性负载均衡
	TrafficMirrorTargetType string `json:"traffic_mirror_target_type"`

	// 功能说明：指定VNI，用于区分不同会话的镜像流量 取值范围：0~16777215 默认值：1
	VirtualNetworkId *string `json:"virtual_network_id,omitempty"`

	// 功能说明：最大传输单元MTU 取值范围：1~1460 默认值：96
	PacketLength *string `json:"packet_length,omitempty"`

	// 功能说明：会话优先级 取值范围：1~32766
	Priority int32 `json:"priority"`

	// 功能说明：是否开启会话 取值范围：true、false 默认值：false
	Enabled *string `json:"enabled,omitempty"`

	// 功能说明：支持的镜像源类型 取值范围：     eni：弹性网卡
	Type *string `json:"type,omitempty"`
}

func (o CreateTrafficMirrorSessionOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorSessionOption struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorSessionOption", string(data)}, " ")
}
