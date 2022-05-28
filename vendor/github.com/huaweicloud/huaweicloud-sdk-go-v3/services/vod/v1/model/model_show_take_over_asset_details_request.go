package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowTakeOverAssetDetailsRequest struct {

	// 媒资原始输入存放的桶。
	SourceBucket string `json:"source_bucket"`

	// 媒资原始输入的objectKey。
	SourceObject string `json:"source_object"`
}

func (o ShowTakeOverAssetDetailsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTakeOverAssetDetailsRequest struct{}"
	}

	return strings.Join([]string{"ShowTakeOverAssetDetailsRequest", string(data)}, " ")
}
