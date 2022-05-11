package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type MigrateSubNetworkInterfaceOption struct {

	// 目的宿主网卡ID
	ParentId string `json:"parent_id"`

	// 待迁移辅助弹性网卡列表
	SubNetworkInterfaces []map[string]string `json:"sub_network_interfaces"`
}

func (o MigrateSubNetworkInterfaceOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MigrateSubNetworkInterfaceOption struct{}"
	}

	return strings.Join([]string{"MigrateSubNetworkInterfaceOption", string(data)}, " ")
}
