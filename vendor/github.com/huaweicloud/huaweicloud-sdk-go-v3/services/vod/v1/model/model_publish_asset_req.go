package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PublishAssetReq struct {

	// 媒资ID。
	AssetId []string `json:"asset_id"`
}

func (o PublishAssetReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PublishAssetReq struct{}"
	}

	return strings.Join([]string{"PublishAssetReq", string(data)}, " ")
}
