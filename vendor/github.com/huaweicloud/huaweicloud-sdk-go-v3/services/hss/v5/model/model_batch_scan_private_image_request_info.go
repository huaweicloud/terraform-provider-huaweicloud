package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchScanPrivateImageRequestInfo struct {

	// 仓库类型，现阶段接入了swr镜像仓库，包含如下:   - SWR : SWR镜像仓库
	RepoType *string `json:"repo_type,omitempty"`

	// 要扫描的镜像信息列表，operate_all参数为false时为必填
	ImageInfoList *[]BatchScanSwrImageInfo `json:"image_info_list,omitempty"`

	// 若为true全量查询，可筛选条件全部查询，若image_info_list为空，则必填
	OperateAll *bool `json:"operate_all,omitempty"`

	// 组织名称
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像版本
	ImageVersion *string `json:"image_version,omitempty"`

	// 镜像类型，包含如下:   - private_image : 私有镜像仓库   - shared_image : 共享镜像仓库
	ImageType string `json:"image_type"`

	// 扫描状态，包含如下:   - unscan : 未扫描   - success : 扫描完成   - scanning : 扫描中   - failed : 扫描失败   - download_failed : 下载失败   - image_oversized : 镜像超大
	ScanStatus *string `json:"scan_status,omitempty"`

	// 仅关注最新版本镜像
	LatestVersion *bool `json:"latest_version,omitempty"`

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
}

func (o BatchScanPrivateImageRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchScanPrivateImageRequestInfo struct{}"
	}

	return strings.Join([]string{"BatchScanPrivateImageRequestInfo", string(data)}, " ")
}
