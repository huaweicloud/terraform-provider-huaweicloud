package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type BatchCreateSubNetworkInterfaceResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	// 批量创建辅助弹性网卡的响应体
	SubNetworkInterfaces *[]SubNetworkInterface `json:"sub_network_interfaces,omitempty"`
	HttpStatusCode       int                    `json:"-"`
}

func (o BatchCreateSubNetworkInterfaceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateSubNetworkInterfaceResponse struct{}"
	}

	return strings.Join([]string{"BatchCreateSubNetworkInterfaceResponse", string(data)}, " ")
}
