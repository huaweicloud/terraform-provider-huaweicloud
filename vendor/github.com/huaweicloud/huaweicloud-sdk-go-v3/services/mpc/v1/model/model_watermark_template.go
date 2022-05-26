package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type WatermarkTemplate struct {

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

	// 图片水印处理方式，type设置为Image时有效。  取值如下：  - Original：只做简单缩放，不做其他处理。 - Grayed：彩色图片变灰。 - Transparent：透明化。
	ImageProcess *string `json:"image_process,omitempty"`

	// 水印图片宽，值有两种形式： - 整数型代水印图片宽的像素值，范围[8，4096]，单位px。 - 小数型代表相对输出视频分辨率宽的比率，范围(0,1)，支持4位小数，如0.9999，超出部分系统自动丢弃。
	Width *string `json:"width,omitempty"`

	// 水印图片高，值有两种形式： - 整数型代表水印图片高的像素值，范围[8，4096]，单位px。 - 小数型代表相对输出视频分辨率高的比率，范围(0，1)，支持4位小数，如0.9999，超出部分系统自动丢弃。
	Height *string `json:"height,omitempty"`

	// 水印叠加母体  取值如下： - input ：水印叠加在输入片源上，转码输出后实际大小按图像等比例缩放 - output ：水印叠加在转码输出文件上。
	Base *WatermarkTemplateBase `json:"base,omitempty"`

	// 水印模板ID
	TemplateId *int32 `json:"template_id,omitempty"`

	// 水印模板名称。
	TemplateName *string `json:"template_name,omitempty"`

	// 水印类型，当前只支持Image（图片水印）。后续根据需求再支持Text（文字水印）。
	Type *string `json:"type,omitempty"`
}

func (o WatermarkTemplate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "WatermarkTemplate struct{}"
	}

	return strings.Join([]string{"WatermarkTemplate", string(data)}, " ")
}

type WatermarkTemplateBase struct {
	value string
}

type WatermarkTemplateBaseEnum struct {
	INPUT  WatermarkTemplateBase
	OUTPUT WatermarkTemplateBase
}

func GetWatermarkTemplateBaseEnum() WatermarkTemplateBaseEnum {
	return WatermarkTemplateBaseEnum{
		INPUT: WatermarkTemplateBase{
			value: "input",
		},
		OUTPUT: WatermarkTemplateBase{
			value: "output",
		},
	}
}

func (c WatermarkTemplateBase) Value() string {
	return c.value
}

func (c WatermarkTemplateBase) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *WatermarkTemplateBase) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
