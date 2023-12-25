package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchScanSwrImageInfo swr私有镜像信息，批量查询需要的参数
type BatchScanSwrImageInfo struct {

	// 命名空间
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像版本
	ImageVersion *string `json:"image_version,omitempty"`
}

func (o BatchScanSwrImageInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchScanSwrImageInfo struct{}"
	}

	return strings.Join([]string{"BatchScanSwrImageInfo", string(data)}, " ")
}
