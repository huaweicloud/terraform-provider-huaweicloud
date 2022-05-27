package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MetaData struct {

	// 文件大小。
	Size *int64 `json:"size,omitempty"`

	// 视频时长，带小数位显示。单位：秒。
	DurationMs *float64 `json:"duration_ms,omitempty"`

	// 视频时长。单位：秒。
	Duration *int64 `json:"duration,omitempty"`

	// 文件封装格式。
	Format *string `json:"format,omitempty"`

	// 总码率。单位：bit/秒
	Bitrate *int64 `json:"bitrate,omitempty"`

	// 视频流元数据。
	Video *[]VideoInfo `json:"video,omitempty"`

	// 音频流元数据。
	Audio *[]AudioInfo `json:"audio,omitempty"`
}

func (o MetaData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MetaData struct{}"
	}

	return strings.Join([]string{"MetaData", string(data)}, " ")
}
