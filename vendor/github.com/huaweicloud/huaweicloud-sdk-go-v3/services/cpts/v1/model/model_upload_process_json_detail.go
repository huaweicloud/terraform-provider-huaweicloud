package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UploadProcessJsonDetail struct {

	// 导入进度Id
	Id *int32 `json:"id,omitempty"`

	// 工程名称
	Name *string `json:"name,omitempty"`

	// 导入状态（0：导入中；1：成功；2：失败）
	Status *int32 `json:"status,omitempty"`

	// 失败原因
	Cause *string `json:"cause,omitempty"`
}

func (o UploadProcessJsonDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UploadProcessJsonDetail struct{}"
	}

	return strings.Join([]string{"UploadProcessJsonDetail", string(data)}, " ")
}
