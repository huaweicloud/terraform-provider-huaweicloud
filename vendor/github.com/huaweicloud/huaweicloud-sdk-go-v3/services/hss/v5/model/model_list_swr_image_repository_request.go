package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSwrImageRepositoryRequest Request Object
type ListSwrImageRepositoryRequest struct {

	// region id
	Region string `json:"region"`

	// 租户企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 组织名称
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称 id
	ImageName *string `json:"image_name,omitempty"`

	// 镜像版本
	ImageVersion *string `json:"image_version,omitempty"`

	// 仅关注最新版本镜像
	LatestVersion *bool `json:"latest_version,omitempty"`

	// 偏移量：指定返回记录的开始位置，必须为数字，取值范围为大于或等于0，默认0
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 镜像类型，包含如下:   - private_image : 私有镜像仓库   - shared_image : 共享镜像仓库
	ImageType string `json:"image_type"`

	// 扫描状态，包含如下:   - unscan : 未扫描   - success : 扫描完成   - scanning : 扫描中   - failed : 扫描失败   - download_failed : 下载失败   - image_oversized : 镜像超大
	ScanStatus *string `json:"scan_status,omitempty"`
}

func (o ListSwrImageRepositoryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSwrImageRepositoryRequest struct{}"
	}

	return strings.Join([]string{"ListSwrImageRepositoryRequest", string(data)}, " ")
}
