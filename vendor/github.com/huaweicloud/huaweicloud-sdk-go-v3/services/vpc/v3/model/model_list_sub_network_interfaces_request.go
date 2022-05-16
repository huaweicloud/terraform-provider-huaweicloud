package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListSubNetworkInterfacesRequest struct {

	// 功能说明：每页返回的个数 取值范围：0-2000
	Limit *int32 `json:"limit,omitempty"`

	// 分页查询起始的资源ID，为空时查询第一页
	Marker *string `json:"marker,omitempty"`

	// 功能说明：辅助弹性网卡ID，支持多ID过滤 使用场景：查询需要的多个辅助弹性网卡信息
	Id *[]string `json:"id,omitempty"`

	// 功能说明：辅助弹性网卡所属虚拟子网的ID，支持多个ID过滤 使用场景：过滤需要的单个或者多个虚拟子网下的辅助弹性网卡
	VirsubnetId *[]string `json:"virsubnet_id,omitempty"`

	// 功能说明：辅助弹性网卡的私有IPv4地址，支持多个地址同时过滤 使用场景：通过单个或者多个ip地址过滤查询辅助弹性网卡
	PrivateIpAddress *[]string `json:"private_ip_address,omitempty"`

	// 功能说明：辅助弹性网卡的mac地址，支持多个同时过滤 使用场景：使用mac地址精确过滤辅助弹性网卡
	MacAddress *[]string `json:"mac_address,omitempty"`

	// 功能说明：辅助弹性网卡所属的VPC_ID，支持多ID过滤 使用场景：过滤单个或多个VPC下的辅助弹性网卡信息
	VpcId *[]string `json:"vpc_id,omitempty"`

	// 功能说明：辅助弹性网卡的描述信息，支持多个同时过滤 使用场景：通过描述信息过滤辅助弹性网卡
	Description *[]string `json:"description,omitempty"`

	// 功能说明：辅助弹性网卡的宿主网卡的ID，支持多ID过滤 使用场景：过滤单个或多个宿主网卡下存在的辅助弹性网卡
	ParentId *[]string `json:"parent_id,omitempty"`
}

func (o ListSubNetworkInterfacesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSubNetworkInterfacesRequest struct{}"
	}

	return strings.Join([]string{"ListSubNetworkInterfacesRequest", string(data)}, " ")
}
