package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DedicatedResourceCapacity struct {
	// 内存大小，单位GB

	Ram *int32 `json:"ram,omitempty"`
	// 磁盘容量，单位GB

	Volume *int64 `json:"volume,omitempty"`
	// cpu核数

	Vcpus *int32 `json:"vcpus,omitempty"`
}

func (o DedicatedResourceCapacity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DedicatedResourceCapacity struct{}"
	}

	return strings.Join([]string{"DedicatedResourceCapacity", string(data)}, " ")
}
