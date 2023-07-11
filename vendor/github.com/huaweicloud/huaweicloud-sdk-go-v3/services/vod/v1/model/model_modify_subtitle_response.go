package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifySubtitleResponse Response Object
type ModifySubtitleResponse struct {

	// 媒资ID。
	AssetId        *string `json:"asset_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ModifySubtitleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifySubtitleResponse struct{}"
	}

	return strings.Join([]string{"ModifySubtitleResponse", string(data)}, " ")
}
