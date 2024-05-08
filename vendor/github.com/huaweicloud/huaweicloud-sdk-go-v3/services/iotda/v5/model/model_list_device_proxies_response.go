package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDeviceProxiesResponse Response Object
type ListDeviceProxiesResponse struct {

	// 代理设备列表
	DeviceProxies *[]QueryDeviceProxySimplify `json:"device_proxies,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ListDeviceProxiesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDeviceProxiesResponse struct{}"
	}

	return strings.Join([]string{"ListDeviceProxiesResponse", string(data)}, " ")
}
