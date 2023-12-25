package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EniNetworkUpdate struct {

	// IPv4子网ID列表。1.19.10及以上版本的CCE Turbo集群支持多容器子网，同时支持增量更新容器子网列表。 只允许新增子网，不可删除已有子网，请谨慎选择。  请求体中需包含所有已经存在的subnet。
	Subnets *[]NetworkSubnet `json:"subnets,omitempty"`
}

func (o EniNetworkUpdate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EniNetworkUpdate struct{}"
	}

	return strings.Join([]string{"EniNetworkUpdate", string(data)}, " ")
}
