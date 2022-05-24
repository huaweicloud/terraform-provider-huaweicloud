package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 媒资审核参数
type Review struct {

	// 审核模板ID。您可以在视频点播控制台配置审核模板后获取，具体请参见[审核设置](https://support.huaweicloud.com/usermanual-vod/vod_01_0057.html)。
	TemplateId string `json:"template_id"`
}

func (o Review) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Review struct{}"
	}

	return strings.Join([]string{"Review", string(data)}, " ")
}
