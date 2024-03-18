package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UploadProcessJson json
type UploadProcessJson struct {

	// 工程导入进度明细信息
	Details *[]UploadProcessJsonDetail `json:"details,omitempty"`

	// 总状态（0：导入中；1：导入完成）
	ProcessStatus *int32 `json:"process_status,omitempty"`
}

func (o UploadProcessJson) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UploadProcessJson struct{}"
	}

	return strings.Join([]string{"UploadProcessJson", string(data)}, " ")
}
