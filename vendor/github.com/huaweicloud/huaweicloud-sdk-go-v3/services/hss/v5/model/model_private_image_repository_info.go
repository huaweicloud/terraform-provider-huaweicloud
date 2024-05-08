package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PrivateImageRepositoryInfo repository info
type PrivateImageRepositoryInfo struct {

	// id
	Id *int64 `json:"id,omitempty"`

	// 命名空间
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像id
	ImageId *string `json:"image_id,omitempty"`

	// 镜像digest
	ImageDigest *string `json:"image_digest,omitempty"`

	// 镜像版本
	ImageVersion *string `json:"image_version,omitempty"`

	// 镜像类型，包含如下2种。   - private_image ：私有镜像。   - shared_image ：共享镜像。
	ImageType *string `json:"image_type,omitempty"`

	// 是否是最新版本
	LatestVersion *bool `json:"latest_version,omitempty"`

	// 扫描状态，包含如下2种。   - unscan ：未扫描。   - success ：扫描完成。   - scanning ：正在扫描。   - failed ：扫描失败。   - download_failed ：下载失败。   - image_oversized ：镜像超大。   - waiting_for_scan ：等待扫描。
	ScanStatus *string `json:"scan_status,omitempty"`

	// 扫描失败原因，包含如下14种。   - \"unknown_error\" :未知错误   - \"authentication_failed\":认证失败   - \"download_failed\":镜像下载失败   - \"image_over_sized\":镜像大小超限   - \"image_oversized\":镜像超大   - \"failed_to_scan_vulnerability\":漏洞扫描失败      - \"failed_to_scan_file\":文件扫描失败   - \"failed_to_scan_software\":软件扫描失败   - \"failed_to_check_sensitive_information\":敏感信息核查失败   - \"failed_to_check_baseline\":基线检查失败   - \"failed_to_check_software_compliance\":软件合规检查失败   - \"failed_to_query_basic_image_information\":基础镜像信息查询失败   - \"response_timed_out\":响应超时   - \"database_error\" : 数据库错误   - \"failed_to_send_the_scan_request\" : 发送扫描请求失败
	ScanFailedDesc *string `json:"scan_failed_desc,omitempty"`

	// 镜像大小
	ImageSize *int64 `json:"image_size,omitempty"`

	// 镜像版本最后更新时间，时间单位 毫秒（ms）
	LatestUpdateTime *int64 `json:"latest_update_time,omitempty"`

	// 最近扫描时间，时间单位 毫秒（ms）
	LatestScanTime *int64 `json:"latest_scan_time,omitempty"`

	// 漏洞个数
	VulNum *int32 `json:"vul_num,omitempty"`

	// 基线扫描未通过数
	UnsafeSettingNum *int32 `json:"unsafe_setting_num,omitempty"`

	// 恶意文件数
	MaliciousFileNum *int32 `json:"malicious_file_num,omitempty"`

	// 拥有者（共享镜像参数）
	DomainName *string `json:"domain_name,omitempty"`

	// 共享镜像状态，包含如下2种。   - expired ：已过期。   - effective ：有效。
	SharedStatus *string `json:"shared_status,omitempty"`

	// 是否可扫描
	Scannable *bool `json:"scannable,omitempty"`

	// 企业版镜像实例名称
	InstanceName *string `json:"instance_name,omitempty"`

	// 企业版镜像实例ID
	InstanceId *string `json:"instance_id,omitempty"`

	// 企业版镜像实例URL
	InstanceUrl *string `json:"instance_url,omitempty"`

	// 多架构关联镜像信息
	AssociationImages *[]AssociateImages `json:"association_images,omitempty"`
}

func (o PrivateImageRepositoryInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrivateImageRepositoryInfo struct{}"
	}

	return strings.Join([]string{"PrivateImageRepositoryInfo", string(data)}, " ")
}
