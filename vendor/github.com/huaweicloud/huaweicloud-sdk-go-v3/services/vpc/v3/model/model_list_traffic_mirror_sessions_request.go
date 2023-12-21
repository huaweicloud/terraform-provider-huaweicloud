package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTrafficMirrorSessionsRequest Request Object
type ListTrafficMirrorSessionsRequest struct {

	// 使用镜像会话ID过滤或排序
	Id *string `json:"id,omitempty"`

	// 使用镜像会话名称过滤或排序
	Name *string `json:"name,omitempty"`

	// 使用镜像会话描述过滤
	Description *string `json:"description,omitempty"`

	// 使用筛选条件ID过滤
	TrafficMirrorFilterId *string `json:"traffic_mirror_filter_id,omitempty"`

	// 使用镜像目的ID过滤
	TrafficMirrorTargetId *string `json:"traffic_mirror_target_id,omitempty"`

	// 使用镜像目的类型过滤
	TrafficMirrorTargetType *string `json:"traffic_mirror_target_type,omitempty"`

	// 使用VNI过滤
	VirtualNetworkId *string `json:"virtual_network_id,omitempty"`

	// 使用最大传输单元MTU过滤
	PacketLength *string `json:"packet_length,omitempty"`

	// 使用镜像会话优先级过滤
	Priority *string `json:"priority,omitempty"`

	// 使用enabled过滤
	Enabled *string `json:"enabled,omitempty"`

	// 使用镜像源类型过滤
	Type *string `json:"type,omitempty"`

	// 使用创建时间戳排序
	CreatedAt *string `json:"created_at,omitempty"`

	// 使用更新时间戳排序
	UpdatedAt *string `json:"updated_at,omitempty"`

	// 功能说明：每页返回的个数 取值范围：0-2000
	Limit *int32 `json:"limit,omitempty"`

	// 分页查询起始的资源ID，为空时查询第一页
	Marker *string `json:"marker,omitempty"`
}

func (o ListTrafficMirrorSessionsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTrafficMirrorSessionsRequest struct{}"
	}

	return strings.Join([]string{"ListTrafficMirrorSessionsRequest", string(data)}, " ")
}
