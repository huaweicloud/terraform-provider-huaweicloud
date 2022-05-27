package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BasicWatermark struct {

	// 水印图片起点相对输出视频顶点的水平偏移量。  设置方法有如下两种：  - 整数型：表示图片起点水平偏移视频顶点的像素值，单位px。取值范围：[0，4096] - 小数型：表示图片起点相对于视频分辨率宽的水平偏移比率。取值范围：(0，1)，支持4位小数，如0.9999，超出部分系统自动丢弃。  示例：输出视频分辨率宽1920，设置“dx”为“0.1”，“referpos”为“TopRight”（右上角），则水印图片右上角到视频右顶点在水平方向上偏移距离为192。
	Dx *string `json:"dx,omitempty"`

	// 水印图片起点相对输出视频顶点的垂直偏移量。  - 设置方法有如下两种：整数型：表示图片起点垂直偏移视频顶点的像素值，单位px。取值范围：[0，4096] - 小数型：表示图片起点相对于视频分辨率高的垂直偏移比率。取值范围：(0，1)，支持4位小数，如0.9999，超出部分系统自动丢弃。  示例：输出视频分辨率高1080，设置“dy”为“0.1”，“referpos”为“TopRight”（右上角），则水印图片右上角到视频右顶点在垂直方向上的偏移距离为108。
	Dy *string `json:"dy,omitempty"`

	// 水印的位置。  取值如下： - TopRight：右上角。 - TopLeft：左上角。 - BottomRight：右下角。 - BottomLeft：左下角。
	Referpos *string `json:"referpos,omitempty"`

	// 水印开始时间，与“timeline_duration”配合使用。  取值范围：数字。  单位：秒。
	TimelineStart *string `json:"timeline_start,omitempty"`

	// 水印持续时间，与“timeline_start”配合使用。  取值范围：[数字，ToEND]。“ToEND”表示持续到视频结束。  默认值：ToEND。
	TimelineDuration *string `json:"timeline_duration,omitempty"`
}

func (o BasicWatermark) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BasicWatermark struct{}"
	}

	return strings.Join([]string{"BasicWatermark", string(data)}, " ")
}
