package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssociateImages 查询swr镜像仓库镜像列表
type AssociateImages struct {

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像版本
	ImageVersion *string `json:"image_version,omitempty"`

	// 镜像类型
	ImageType *string `json:"image_type,omitempty"`

	// 命名空间
	Namespace *string `json:"namespace,omitempty"`

	// 镜像digest
	ImageDigest *string `json:"image_digest,omitempty"`

	// 扫描状态，包含如下2种。   - unscan ：未扫描。   - success ：扫描完成。   - scanning ：正在扫描。   - failed ：扫描失败。   - download_failed ：下载失败。   - image_oversized ：镜像超大。   - waiting_for_scan ：等待扫描。
	ScanStatus *string `json:"scan_status,omitempty"`
}

func (o AssociateImages) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateImages struct{}"
	}

	return strings.Join([]string{"AssociateImages", string(data)}, " ")
}
