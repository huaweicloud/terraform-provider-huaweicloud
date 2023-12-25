package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MasterEipResponseStatus 状态信息
type MasterEipResponseStatus struct {

	// 集群访问的PrivateIP(HA集群返回VIP)
	PrivateEndpoint *string `json:"privateEndpoint,omitempty"`

	// 集群访问的PublicIP
	PublicEndpoint *string `json:"publicEndpoint,omitempty"`
}

func (o MasterEipResponseStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MasterEipResponseStatus struct{}"
	}

	return strings.Join([]string{"MasterEipResponseStatus", string(data)}, " ")
}
