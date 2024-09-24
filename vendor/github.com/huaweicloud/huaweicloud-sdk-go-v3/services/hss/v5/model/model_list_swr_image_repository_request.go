package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSwrImageRepositoryRequest Request Object
type ListSwrImageRepositoryRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 组织名称
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像版本
	ImageVersion *string `json:"image_version,omitempty"`

	// 仅关注最新版本镜像
	LatestVersion *bool `json:"latest_version,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 镜像类型，包含如下:   - private_image : 私有镜像仓库   - shared_image : 共享镜像仓库   - local_image : 本地镜像   - instance_image : 企业镜像
	ImageType string `json:"image_type"`

	// 扫描状态，包含如下:   - unscan : 未扫描   - success : 扫描完成   - scanning : 扫描中   - failed : 扫描失败   - waiting_for_scan : 等待扫描
	ScanStatus *string `json:"scan_status,omitempty"`

	// 企业镜像实例名称
	InstanceName *string `json:"instance_name,omitempty"`

	// 镜像大小
	ImageSize *int64 `json:"image_size,omitempty"`

	// 创建时间开始日期，时间单位 毫秒（ms）
	StartLatestUpdateTime *int64 `json:"start_latest_update_time,omitempty"`

	// 创建时间结束日期，时间单位 毫秒（ms）
	EndLatestUpdateTime *int64 `json:"end_latest_update_time,omitempty"`

	// 最近一次扫描完成时间开始日期，时间单位 毫秒（ms）
	StartLatestScanTime *int64 `json:"start_latest_scan_time,omitempty"`

	// 最近一次扫描完成时间结束日期，时间单位 毫秒（ms）
	EndLatestScanTime *int64 `json:"end_latest_scan_time,omitempty"`

	// 是否存在恶意文件
	HasMaliciousFile *bool `json:"has_malicious_file,omitempty"`

	// 是否存在基线检查
	HasUnsafeSetting *bool `json:"has_unsafe_setting,omitempty"`

	// 是否存在软件漏洞
	HasVul *bool `json:"has_vul,omitempty"`

	// 企业仓库实例ID，swr共享版无需使用该参数
	InstanceId *string `json:"instance_id,omitempty"`
}

func (o ListSwrImageRepositoryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSwrImageRepositoryRequest struct{}"
	}

	return strings.Join([]string{"ListSwrImageRepositoryRequest", string(data)}, " ")
}
