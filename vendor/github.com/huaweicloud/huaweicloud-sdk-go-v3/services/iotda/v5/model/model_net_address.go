package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 网络地址信息。
type NetAddress struct {

	// **参数说明**：服务的对应IP
	Ip *string `json:"ip,omitempty"`

	// **参数说明**：服务对应端口
	Port *int32 `json:"port,omitempty"`

	// **参数说明**：服务对应的域名
	Domain *string `json:"domain,omitempty"`
}

func (o NetAddress) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NetAddress struct{}"
	}

	return strings.Join([]string{"NetAddress", string(data)}, " ")
}
