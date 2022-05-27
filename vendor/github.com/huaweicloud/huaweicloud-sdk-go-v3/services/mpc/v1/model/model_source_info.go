package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SourceInfo struct {

	// 片源时长，单位：秒
	Duration *int32 `json:"duration,omitempty"`

	// 片源时长，单位：毫秒
	DurationMs *int64 `json:"duration_ms,omitempty"`

	// 片源格式
	Format *string `json:"format,omitempty"`

	// 片源大小
	Size *int64 `json:"size,omitempty"`

	VideoInfo *VideoInfo `json:"video_info,omitempty"`

	// 音频信息
	AudioInfo *[]AudioInfo `json:"audio_info,omitempty"`
}

func (o SourceInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourceInfo struct{}"
	}

	return strings.Join([]string{"SourceInfo", string(data)}, " ")
}
