package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Review 媒资审核参数
type Review struct {

	// 审核模板ID。您可以在视频点播控制台配置审核模板后获取，具体请参见[审核设置](https://support.huaweicloud.com/usermanual-vod/vod_01_0057.html)。
	TemplateId string `json:"template_id"`

	// 截图检测时间间隔，取值范围为[0,100]，该参数在请求参数中忽略。
	Interval *int32 `json:"interval,omitempty"`

	// 鉴政内容检测置信度，取值范围为[0,100]，该参数在请求参数中忽略。 置信度越高，说明审核结果越可信。未开启或设置为0时，表示未进行此项检测。
	Politics *int32 `json:"politics,omitempty"`

	// 鉴恐内容的检测置信度，取值范围为[0,100]，该参数在请求参数中忽略。 置信度越高，说明审核结果越可信。未开启或设置为0时，表示未进行此项检测。
	Terrorism *int32 `json:"terrorism,omitempty"`

	// 鉴黄内容的检测置信度，取值范围为[0,100]，该参数在请求参数中忽略。 置信度越高，说明审核结果越可信。未开启或设置为0时，表示未进行此项检测。
	Porn *int32 `json:"porn,omitempty"`
}

func (o Review) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Review struct{}"
	}

	return strings.Join([]string{"Review", string(data)}, " ")
}
