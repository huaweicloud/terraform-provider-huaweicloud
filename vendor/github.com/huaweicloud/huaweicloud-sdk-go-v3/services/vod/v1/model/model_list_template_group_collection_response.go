package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTemplateGroupCollectionResponse struct {

	// 模板组集合信息<br/>
	TemplateGroupCollectionList *[]TemplateGroupCollection `json:"template_group_collection_list,omitempty"`

	// 总记录条数<br/>
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListTemplateGroupCollectionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTemplateGroupCollectionResponse struct{}"
	}

	return strings.Join([]string{"ListTemplateGroupCollectionResponse", string(data)}, " ")
}
