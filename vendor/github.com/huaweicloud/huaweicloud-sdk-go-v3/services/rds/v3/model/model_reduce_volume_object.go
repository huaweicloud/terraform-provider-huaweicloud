package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ReduceVolumeObject 实例磁盘缩容时必填。
type ReduceVolumeObject struct {

	// 缩容后实例磁盘的目标大小。每次缩容至少缩小10GB；目标大小必须为10的整数倍。 为确保实例的正常使用，根据当前磁盘的使用量情况存在磁盘容量下限，当此参数小于磁盘容量下限时，缩容会下发失败，此时请适当调大此参数。
	Size int32 `json:"size"`
}

func (o ReduceVolumeObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReduceVolumeObject struct{}"
	}

	return strings.Join([]string{"ReduceVolumeObject", string(data)}, " ")
}
