package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type WatermarkTemplate struct {

	// 水印模板名称。
	Name *string `json:"name,omitempty"`

	// 水印模板配置id。
	Id *string `json:"id,omitempty"`

	// 启用状态。  取值为： - 0：停用 - 1：启用
	Status *int64 `json:"status,omitempty"`

	// 水印图片相对输出视频的水平偏移量。  默认值是0。
	Dx *string `json:"dx,omitempty"`

	// 水印图片相对输出视频的垂直偏移量。  默认值是0。
	Dy *string `json:"dy,omitempty"`

	// 水印的位置。
	Position *string `json:"position,omitempty"`

	// 水印图片宽。
	Width *string `json:"width,omitempty"`

	// 水印图片高。
	Height *string `json:"height,omitempty"`

	// 创建时间。
	CreateTime *string `json:"create_time,omitempty"`

	// 水印图片下载url。
	ImageUrl *string `json:"image_url,omitempty"`

	// 水印图片格式类型。
	Type *string `json:"type,omitempty"`

	// 水印类型，当前只支持Image（图片水印）。
	WatermarkType *string `json:"watermark_type,omitempty"`

	// type设置为Image时有效。  目前包括： - Original：只做简单缩放，不做其他处理 - Transparent：图片底色透明 - Grayed：彩色图片变灰
	ImageProcess *string `json:"image_process,omitempty"`

	// 水印开始时间，与\"timeline_duration\"配合使用。 取值范围:[0, END)。\"END\"表示视频结束时间。 单位:秒。
	TimelineStart *string `json:"timeline_start,omitempty"`

	// 水印持续时间，与\"timeline_start\"配合使用。 取值范围:(0,END-开始时间]。\"END\"表示视频结束时间。 单位:秒。 默认:END。
	TimelineDuration *string `json:"timeline_duration,omitempty"`
}

func (o WatermarkTemplate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "WatermarkTemplate struct{}"
	}

	return strings.Join([]string{"WatermarkTemplate", string(data)}, " ")
}
