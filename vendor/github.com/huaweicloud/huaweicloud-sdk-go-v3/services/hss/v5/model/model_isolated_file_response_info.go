package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IsolatedFileResponseInfo 已隔离文件详情
type IsolatedFileResponseInfo struct {

	// 操作系统类型，包含如下2种。   - Linux ：Linux。   - Windows ：Windows。
	OsType string `json:"os_type"`

	// 主机ID
	HostId string `json:"host_id"`

	// 服务器名称
	HostName string `json:"host_name"`

	// 文件哈希
	FileHash string `json:"file_hash"`

	// 文件路径
	FilePath string `json:"file_path"`

	// 文件属性
	FileAttr string `json:"file_attr"`

	// 隔离状态，包含如下:   - isolated : 已隔离   - restored : 已恢复   - isolating : 已下发隔离任务   - restoring : 已下发恢复任务
	IsolationStatus string `json:"isolation_status"`

	// 服务器私有IP
	PrivateIp string `json:"private_ip"`

	// 弹性公网IP地址
	PublicIp string `json:"public_ip"`

	// 资产重要性
	AssetValue string `json:"asset_value"`

	// 更新时间，毫秒
	UpdateTime int64 `json:"update_time"`

	// agent版本
	AgentVersion string `json:"agent_version"`

	// 隔离来源，包含如下:   - event : 安全告警事件   - antivirus : 病毒查杀
	IsolateSource string `json:"isolate_source"`

	// 事件名称
	EventName string `json:"event_name"`

	AgentEventInfo *IsolateEventResponseInfo `json:"agent_event_info"`

	AntivirusResultInfo *AntivirusResultDetailInfo `json:"antivirus_result_info"`
}

func (o IsolatedFileResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IsolatedFileResponseInfo struct{}"
	}

	return strings.Join([]string{"IsolatedFileResponseInfo", string(data)}, " ")
}
