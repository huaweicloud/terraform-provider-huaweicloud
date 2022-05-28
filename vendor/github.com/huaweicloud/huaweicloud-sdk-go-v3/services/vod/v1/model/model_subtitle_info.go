package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SubtitleInfo struct {

	// 字幕文件的下载地址
	Url *string `json:"url,omitempty"`

	// 字幕文件id
	Id *int32 `json:"id,omitempty"`

	// 字幕文件类型
	Type *string `json:"type,omitempty"`

	// 字幕文件语言种类
	Language *string `json:"language,omitempty"`
}

func (o SubtitleInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SubtitleInfo struct{}"
	}

	return strings.Join([]string{"SubtitleInfo", string(data)}, " ")
}
