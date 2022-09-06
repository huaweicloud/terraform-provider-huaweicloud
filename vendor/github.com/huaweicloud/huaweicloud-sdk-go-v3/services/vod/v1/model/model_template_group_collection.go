package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TemplateGroupCollection struct {

	// 模板组集合id<br/>
	GroupCollectionId *string `json:"group_collection_id,omitempty"`

	// 模板组集合名称<br/>
	Name *string `json:"name,omitempty"`

	// 模板介绍<br/>
	Description *string `json:"description,omitempty"`

	// 转码组列表<br/>
	TemplateGroupList *[]TemplateGroup `json:"template_group_list,omitempty"`
}

func (o TemplateGroupCollection) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TemplateGroupCollection struct{}"
	}

	return strings.Join([]string{"TemplateGroupCollection", string(data)}, " ")
}
