package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TracksInfo struct {

	// 音频轨的声道layout
	ChannelLayout *string `json:"channel_layout,omitempty"`

	// 音频轨对应语言描述
	Language *string `json:"language,omitempty"`
}

func (o TracksInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TracksInfo struct{}"
	}

	return strings.Join([]string{"TracksInfo", string(data)}, " ")
}
