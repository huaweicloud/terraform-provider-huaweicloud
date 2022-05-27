package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type WatermarkRequest struct {
	Input *ObsObjInfo `json:"input,omitempty"`

	// 水印模板。可通过新建水印模板接口创建水印模板。
	TemplateId *string `json:"template_id,omitempty"`

	// 文字水印内容，内容需做Base64编码，若类型为文字水印 (type字段为Text)，则此配置项不能为空  示例：若想添加文字水印“测试文字水印”，那么Content的值为：5rWL6K+V5paH5a2X5rC05Y2w
	TextContext *string `json:"text_context,omitempty"`

	ImageWatermark *ImageWatermark `json:"image_watermark,omitempty"`

	TextWatermark *TextWatermark `json:"text_watermark,omitempty"`
}

func (o WatermarkRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "WatermarkRequest struct{}"
	}

	return strings.Join([]string{"WatermarkRequest", string(data)}, " ")
}
