package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeBandwidth
type NodeBandwidth struct {

	// 带宽的计费类型： - 未传该字段，表示按带宽计费。 - 字段值为空，表示按带宽计费。 - 字段值为“traffic”，表示按流量计费。 - 字段为其它值，会导致创建云服务器失败。 > - 按带宽计费：按公网传输速率（单位为Mbps）计费。当您的带宽利用率高于10%时，建议优先选择按带宽计费。 > - 按流量计费：只允许在创建按需节点时指定，按公网传输的数据总量（单位为GB）计费。当您的带宽利用率低于10%时，建议优先选择按流量计费。
	Chargemode *string `json:"chargemode,omitempty"`

	// 带宽大小，取值请参见取值请参见申请EIP接口中bandwidth.size说明。 [链接请参见[申请EIP](https://support.huaweicloud.com/api-eip/eip_api_0001.html)](tag:hws) [链接请参见[申请EIP](https://support.huaweicloud.com/intl/zh-cn/api-eip/eip_api_0001.html)](tag:hws_hk)
	Size *int32 `json:"size,omitempty"`

	// 带宽的共享类型，共享类型枚举：PER，表示独享，目前仅支持独享。
	Sharetype *string `json:"sharetype,omitempty"`
}

func (o NodeBandwidth) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeBandwidth struct{}"
	}

	return strings.Join([]string{"NodeBandwidth", string(data)}, " ")
}
