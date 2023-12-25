package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RunImageSynchronizeRequestInfo struct {

	// 镜像类型，包含如下:   - private_image : 私有镜像仓库   - shared_image : 共享镜像仓库
	ImageType string `json:"image_type"`
}

func (o RunImageSynchronizeRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RunImageSynchronizeRequestInfo struct{}"
	}

	return strings.Join([]string{"RunImageSynchronizeRequestInfo", string(data)}, " ")
}
