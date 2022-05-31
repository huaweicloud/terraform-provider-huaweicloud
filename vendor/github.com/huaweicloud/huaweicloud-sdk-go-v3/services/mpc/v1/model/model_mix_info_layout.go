package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MixInfoLayout struct {

	// 原视频在合成视频中的位置布局配置
	Panes []PaneSetting `json:"panes"`
}

func (o MixInfoLayout) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MixInfoLayout struct{}"
	}

	return strings.Join([]string{"MixInfoLayout", string(data)}, " ")
}
