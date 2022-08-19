package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TransTemplateGroupCollection struct {

	// 模板组集合名称<br/>
	Name *string `json:"name,omitempty"`

	// 模板组集合描述<br/>
	Description *string `json:"description,omitempty"`

	// 模板组列表,模板ID<br/>
	TemplateGroupList *[]string `json:"template_group_list,omitempty"`
}

func (o TransTemplateGroupCollection) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TransTemplateGroupCollection struct{}"
	}

	return strings.Join([]string{"TransTemplateGroupCollection", string(data)}, " ")
}
