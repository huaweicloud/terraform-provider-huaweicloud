package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ExtractAudioTaskReq struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`

	Parameter *Parameter `json:"parameter,omitempty"`
}

func (o ExtractAudioTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExtractAudioTaskReq struct{}"
	}

	return strings.Join([]string{"ExtractAudioTaskReq", string(data)}, " ")
}
