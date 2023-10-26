package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeDetailRequest Request Object
type UpgradeDetailRequest struct {

	// 待升级的集群的ID。
	ClusterId string `json:"cluster_id"`

	// 偏移量。 偏移量为一个大于0小于终端节点服务总个数的整数， 表示从偏移量后面的终端节点服务开始查询。
	Start *int32 `json:"start,omitempty"`

	// 查询返回终端节点服务的连接列表限制每页个数，即每页返回的个数。
	Limit *int32 `json:"limit,omitempty"`

	// 查询升级行为。 - 查询集群版本升级详情：不填写该参数。 - 查询切换AZ详情：当前仅支持AZ_MIGRATION。
	ActionMode *string `json:"action_mode,omitempty"`
}

func (o UpgradeDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDetailRequest struct{}"
	}

	return strings.Join([]string{"UpgradeDetailRequest", string(data)}, " ")
}
