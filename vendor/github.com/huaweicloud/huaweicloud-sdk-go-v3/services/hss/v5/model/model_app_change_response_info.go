package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AppChangeResponseInfo 软件变动历史信息
type AppChangeResponseInfo struct {

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// the type of change   - add ：新建   - delete ：删除   - modify ：修改
	VariationType *string `json:"variation_type,omitempty"`

	// host_id
	HostId *string `json:"host_id,omitempty"`

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器ip
	HostIp *string `json:"host_ip,omitempty"`

	// 版本号
	Version *string `json:"version,omitempty"`

	// 软件更新时间
	UpdateTime *int64 `json:"update_time,omitempty"`

	// 最近扫描时间
	RecentScanTime *int64 `json:"recent_scan_time,omitempty"`
}

func (o AppChangeResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AppChangeResponseInfo struct{}"
	}

	return strings.Join([]string{"AppChangeResponseInfo", string(data)}, " ")
}
