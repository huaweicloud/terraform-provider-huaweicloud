package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VideoAndTemplate struct {

	// 模板ID
	TemplateId *int32 `json:"template_id,omitempty"`

	// 视频宽度（单位：像素） - H264：范围[32,4096]，必须为2的倍数 - H265：范围[320,4096]，必须是4的倍数
	Width *int32 `json:"width,omitempty"`

	// 视频高度（单位：像素） - H264：范围[32,2880]，必须为2的倍数 - H265：范围[240,2880]，必须是4的倍数
	Height *int32 `json:"height,omitempty"`

	// 输出平均码率。  取值范围：0或[40,30000]之间的整数。  单位：kbit/s  若设置为0，则输出平均码率为自适应值。
	Bitrate *int32 `json:"bitrate,omitempty"`
}

func (o VideoAndTemplate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoAndTemplate struct{}"
	}

	return strings.Join([]string{"VideoAndTemplate", string(data)}, " ")
}
