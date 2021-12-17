package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 网卡资源
type Nics struct {
	// 子网ID，如果是自动创建，使用\"autoCreate\"

	Id string `json:"id"`
	// 子网名称

	Name string `json:"name"`
	// 子网网关/掩码

	Cidr string `json:"cidr"`
	// 虚拟机IP地址，如果没有这个字段，自动分配IP

	Ip *string `json:"ip,omitempty"`
}

func (o Nics) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Nics struct{}"
	}

	return strings.Join([]string{"Nics", string(data)}, " ")
}
