package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CancelExtractAudioTaskRequest struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`
}

func (o CancelExtractAudioTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CancelExtractAudioTaskRequest struct{}"
	}

	return strings.Join([]string{"CancelExtractAudioTaskRequest", string(data)}, " ")
}
