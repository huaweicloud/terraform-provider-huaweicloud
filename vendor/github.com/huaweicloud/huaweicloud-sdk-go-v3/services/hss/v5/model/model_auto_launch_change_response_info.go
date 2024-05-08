package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutoLaunchChangeResponseInfo 自启动项变动历史信息
type AutoLaunchChangeResponseInfo struct {

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// the type of change   - add ：新建   - delete ：删除   - modify ：修改
	VariationType *string `json:"variation_type,omitempty"`

	// 自启动项类型   - 0 ：自启动服务   - 1 ：定时任务   - 2 ：预加载动态库   - 3 ：Run注册表键   - 4 ：开机启动文件夹
	Type *int32 `json:"type,omitempty"`

	// host_id
	HostId *string `json:"host_id,omitempty"`

	// 弹性服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 主机IP
	HostIp *string `json:"host_ip,omitempty"`

	// 自启动项的路径
	Path *string `json:"path,omitempty"`

	// 采用sha256算法生成的文件hash值
	Hash *string `json:"hash,omitempty"`

	// 运行用户
	RunUser *string `json:"run_user,omitempty"`

	// 自启动项名称
	Name *string `json:"name,omitempty"`

	// 最近更新时间，13位时间戳
	RecentScanTime *int64 `json:"recent_scan_time,omitempty"`
}

func (o AutoLaunchChangeResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutoLaunchChangeResponseInfo struct{}"
	}

	return strings.Join([]string{"AutoLaunchChangeResponseInfo", string(data)}, " ")
}
