package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VodInfoV2 struct {

	// VOD媒资id
	AssetId string `json:"asset_id"`

	// 点播播放地址
	PlayUrl *string `json:"play_url,omitempty"`
}

func (o VodInfoV2) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VodInfoV2 struct{}"
	}

	return strings.Join([]string{"VodInfoV2", string(data)}, " ")
}
