package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AutopilotEniNetworkUpdate struct {

	// IPv4子网ID列表。 只允许新增子网，不可删除已有子网，请谨慎选择。  请求体中需包含所有已经存在的subnet。
	Subnets *[]AutopilotNetworkSubnet `json:"subnets,omitempty"`
}

func (o AutopilotEniNetworkUpdate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotEniNetworkUpdate struct{}"
	}

	return strings.Join([]string{"AutopilotEniNetworkUpdate", string(data)}, " ")
}
