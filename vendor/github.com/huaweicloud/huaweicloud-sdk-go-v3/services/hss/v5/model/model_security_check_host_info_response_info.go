package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SecurityCheckHostInfoResponseInfo 受配置检测影响的服务器信息
type SecurityCheckHostInfoResponseInfo struct {

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器公网IP
	HostPublicIp *string `json:"host_public_ip,omitempty"`

	// 服务器私网IP
	HostPrivateIp *string `json:"host_private_ip,omitempty"`

	// 扫描时间(ms)
	ScanTime *int64 `json:"scan_time,omitempty"`

	// 风险项数量
	FailedNum *int32 `json:"failed_num,omitempty"`

	// 通过项数量
	PassedNum *int32 `json:"passed_num,omitempty"`
}

func (o SecurityCheckHostInfoResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SecurityCheckHostInfoResponseInfo struct{}"
	}

	return strings.Join([]string{"SecurityCheckHostInfoResponseInfo", string(data)}, " ")
}
