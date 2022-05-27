package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type OutputVideoPara struct {

	// 输出视频对应的模板ID
	TemplateId *int32 `json:"template_id,omitempty"`

	// 视频大小
	Size *int64 `json:"size,omitempty"`

	// 视频封装格式
	Pack *string `json:"pack,omitempty"`

	Video *VideoInfo `json:"video,omitempty"`

	Audio *AudioInfo `json:"audio,omitempty"`

	// 输出片源文件名
	FileName *string `json:"file_name,omitempty"`

	// 折算后视频时长
	ConverDuration *float64 `json:"conver_duration,omitempty"`

	Error *XCodeError `json:"error,omitempty"`
}

func (o OutputVideoPara) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OutputVideoPara struct{}"
	}

	return strings.Join([]string{"OutputVideoPara", string(data)}, " ")
}
