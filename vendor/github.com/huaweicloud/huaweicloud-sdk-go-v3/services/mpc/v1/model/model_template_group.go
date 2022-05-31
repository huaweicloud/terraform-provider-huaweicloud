package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TemplateGroup struct {

	// 模板组id
	GroupId *string `json:"group_id,omitempty"`

	// 模板组名称
	Name *string `json:"name,omitempty"`

	// 模板组模板ID
	TemplateIds *[]int32 `json:"template_ids,omitempty"`

	// 视频信息列表
	Videos *[]VideoAndTemplate `json:"videos,omitempty"`

	Audio *Audio `json:"audio,omitempty"`

	VideoCommon *VideoCommon `json:"video_common,omitempty"`

	Common *Common `json:"common,omitempty"`
}

func (o TemplateGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TemplateGroup struct{}"
	}

	return strings.Join([]string{"TemplateGroup", string(data)}, " ")
}
