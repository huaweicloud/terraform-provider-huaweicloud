package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateExtractAudioTaskResponse struct {

	// 视频源媒资ID。
	AssetId *string `json:"asset_id,omitempty"`

	// 提取的音频媒资ID。
	AudioAssetId   *string `json:"audio_asset_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateExtractAudioTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateExtractAudioTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateExtractAudioTaskResponse", string(data)}, " ")
}
