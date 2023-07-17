package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterVolumeRsp 实例磁盘信息。
type ShowClusterVolumeRsp struct {

	// 实例磁盘类型。
	Type *string `json:"type,omitempty"`

	// 实例磁盘大小。
	Size *int32 `json:"size,omitempty"`
}

func (o ShowClusterVolumeRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterVolumeRsp struct{}"
	}

	return strings.Join([]string{"ShowClusterVolumeRsp", string(data)}, " ")
}
