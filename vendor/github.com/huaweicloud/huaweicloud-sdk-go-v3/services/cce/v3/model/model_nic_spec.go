package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NicSpec 主网卡的描述信息。
type NicSpec struct {

	// 网卡所在子网的网络ID。主网卡创建时若未指定subnetId,将使用集群子网。若节点池同时配置了subnetList，则节点池扩容子网以subnetList字段为准。扩展网卡创建时必须指定subnetId。
	SubnetId *string `json:"subnetId,omitempty"`

	// 主网卡的IP将通过fixedIps指定，数量不得大于创建的节点数。fixedIps或ipBlock同时只能指定一个。扩展网卡不支持指定fiexdIps。
	FixedIps *[]string `json:"fixedIps,omitempty"`

	// 主网卡的IP段的CIDR格式，创建的节点IP将属于该IP段内。fixedIps或ipBlock同时只能指定一个。
	IpBlock *string `json:"ipBlock,omitempty"`

	// 网卡所在子网的网络ID列表，支持节点池配置多个子网，最多支持配置20个子网。
	SubnetList *[]string `json:"subnetList,omitempty"`
}

func (o NicSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NicSpec struct{}"
	}

	return strings.Join([]string{"NicSpec", string(data)}, " ")
}
