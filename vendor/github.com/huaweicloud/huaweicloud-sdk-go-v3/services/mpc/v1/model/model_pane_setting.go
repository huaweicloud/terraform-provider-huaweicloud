package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PaneSetting struct {

	// 原视频的id。
	PaneId string `json:"pane_id"`

	// pane_id标记的原视频起点，在合成视频中相对于左下角的水平偏移量。 目前只支持小数类型，表示相对于输出视频宽的水平偏移比率。取值范围(0,1)。
	X string `json:"x"`

	// pane_id标记的原视频，在合成视频中相对于左下角的垂直偏移量。 目前只支持小数型，表示相对于输出视频高的垂直偏移比率。取值范围:(0,1)。
	Y string `json:"y"`

	// pane_id标记的原视频，在合成视频中占的宽。目前只支持小数型，范围(0,1)，表示占据合成视频宽的比率。
	Width string `json:"width"`

	// pane_id标记的原视频，在合成视频中占的高。目前只支持小数型，范围(0,1)，表示占据合成视频高的比率。
	Height string `json:"height"`
}

func (o PaneSetting) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PaneSetting struct{}"
	}

	return strings.Join([]string{"PaneSetting", string(data)}, " ")
}
