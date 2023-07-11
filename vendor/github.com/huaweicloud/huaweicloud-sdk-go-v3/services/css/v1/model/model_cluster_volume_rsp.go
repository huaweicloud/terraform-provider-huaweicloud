package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClusterVolumeRsp 实例磁盘信息。
type ClusterVolumeRsp struct {

	// 实例磁盘类型。
	Type *string `json:"type,omitempty"`

	// 实例磁盘大小。
	Size *int32 `json:"size,omitempty"`
}

func (o ClusterVolumeRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterVolumeRsp struct{}"
	}

	return strings.Join([]string{"ClusterVolumeRsp", string(data)}, " ")
}
