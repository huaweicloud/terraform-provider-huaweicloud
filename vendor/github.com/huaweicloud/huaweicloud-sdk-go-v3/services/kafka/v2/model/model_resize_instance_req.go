package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResizeInstanceReq struct {

	// 规格变更后的规格ID。 若只扩展磁盘大小，则规格ID保持和原实例不变。
	NewSpecCode *string `json:"new_spec_code,omitempty"`

	// 规格变更后的消息存储空间，单位：GB。 若扩展实例基准带宽，则new_storage_space不能低于基准带宽规定的最小磁盘大小。
	NewStorageSpace *int32 `json:"new_storage_space,omitempty"`
}

func (o ResizeInstanceReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeInstanceReq struct{}"
	}

	return strings.Join([]string{"ResizeInstanceReq", string(data)}, " ")
}
