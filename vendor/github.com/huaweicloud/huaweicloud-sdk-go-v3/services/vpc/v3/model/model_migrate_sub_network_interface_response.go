package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type MigrateSubNetworkInterfaceResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	// 批量迁移辅助弹性网卡的响应体
	SubNetworkInterfaces *[]SubNetworkInterface `json:"sub_network_interfaces,omitempty"`
	HttpStatusCode       int                    `json:"-"`
}

func (o MigrateSubNetworkInterfaceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MigrateSubNetworkInterfaceResponse struct{}"
	}

	return strings.Join([]string{"MigrateSubNetworkInterfaceResponse", string(data)}, " ")
}
