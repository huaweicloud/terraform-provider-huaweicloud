package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyTransTemplateGroup struct {

	// 模板组ID
	GroupId *string `json:"group_id,omitempty"`

	// 模板组名称
	Name *string `json:"name,omitempty"`

	// 视频信息列表
	Videos *[]VideoObj `json:"videos,omitempty"`

	Audio *Audio `json:"audio,omitempty"`

	VideoCommon *VideoCommon `json:"video_common,omitempty"`

	Common *Common `json:"common,omitempty"`
}

func (o ModifyTransTemplateGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyTransTemplateGroup struct{}"
	}

	return strings.Join([]string{"ModifyTransTemplateGroup", string(data)}, " ")
}
