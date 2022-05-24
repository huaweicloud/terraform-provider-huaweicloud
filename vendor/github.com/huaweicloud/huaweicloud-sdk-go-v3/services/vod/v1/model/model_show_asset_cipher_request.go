package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowAssetCipherRequest struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`
}

func (o ShowAssetCipherRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetCipherRequest struct{}"
	}

	return strings.Join([]string{"ShowAssetCipherRequest", string(data)}, " ")
}
