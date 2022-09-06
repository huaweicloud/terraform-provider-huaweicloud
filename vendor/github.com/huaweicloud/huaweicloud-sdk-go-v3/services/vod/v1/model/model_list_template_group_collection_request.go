package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTemplateGroupCollectionRequest struct {

	// 模板组集合id
	GroupCollectionId *string `json:"group_collection_id,omitempty"`

	// 偏移量。默认为0。指定group_collection_id时该参数无效。<br/>
	Offset *int32 `json:"offset,omitempty"`

	// 每页记录数。默认为10，范围[1,100]。指定group_collection_id时该参数无效。<br/>
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListTemplateGroupCollectionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTemplateGroupCollectionRequest struct{}"
	}

	return strings.Join([]string{"ListTemplateGroupCollectionRequest", string(data)}, " ")
}
