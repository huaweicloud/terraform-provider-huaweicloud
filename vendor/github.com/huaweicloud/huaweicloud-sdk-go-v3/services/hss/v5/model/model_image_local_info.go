package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ImageLocalInfo struct {

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像ID
	ImageId *string `json:"image_id,omitempty"`

	// 镜像digest
	ImageDigest *string `json:"image_digest,omitempty"`

	// 镜像版本
	ImageVersion *string `json:"image_version,omitempty"`

	// 本地镜像类型
	LocalImageType *string `json:"local_image_type,omitempty"`

	// 扫描状态，包含如下：   - unscan：未扫描   - success：扫描完成   - scanning：正在扫描   - failed：扫描失败   - waiting：等待扫描
	ScanStatus *string `json:"scan_status,omitempty"`

	// 镜像大小，单位字节
	ImageSize *int64 `json:"image_size,omitempty"`

	// 镜像版本最后更新时间，时间单位毫秒（ms）
	LatestUpdateTime *int64 `json:"latest_update_time,omitempty"`

	// 最近扫描时间，时间单位毫秒（ms）
	LatestScanTime *int64 `json:"latest_scan_time,omitempty"`

	// 漏洞个数
	VulNum *int64 `json:"vul_num,omitempty"`

	// 基线扫描未通过数
	UnsafeSettingNum *int64 `json:"unsafe_setting_num,omitempty"`

	// 恶意文件数
	MaliciousFileNum *int64 `json:"malicious_file_num,omitempty"`

	// 关联主机数
	HostNum *int64 `json:"host_num,omitempty"`

	// 关联容器数
	ContainerNum *int64 `json:"container_num,omitempty"`

	// 关联组件数
	ComponentNum *int64 `json:"component_num,omitempty"`

	// 扫描失败原因，包含如下10种。   - \"unknown_error\":未知错误   - \"failed_to_match_agent\":对应主机未开启容器版防护或agent离线   - \"create_container_failed\":创建容器失败        - \"get_container_info_failed\":获取容器信息失败   - \"docker_offline\":docker引擎不在线   - \"get_docker_root_failed\":获取容器根文件系统失败   - \"image_not_exist_or_docker_api_fault\":镜像不存在或docker接口错误   - \"huge_image\":超大镜像   - \"docker_root_in_nfs\":容器根目录位于网络挂载   - \"response_timed_out\":响应超时
	ScanFailedDesc *string `json:"scan_failed_desc,omitempty"`

	// 镜像风险程度，在镜像扫描完成后展示，包含如下：   - Security：安全   - Low：低危   - Medium：中危   - High：高危
	SeverityLevel *string `json:"severity_level,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`
}

func (o ImageLocalInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImageLocalInfo struct{}"
	}

	return strings.Join([]string{"ImageLocalInfo", string(data)}, " ")
}
