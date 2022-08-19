package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTemplateGroupCollectionResponse struct {

	// 模板组集合ID<br/>
	GroupCollectionId *string `json:"group_collection_id,omitempty"`
	HttpStatusCode    int     `json:"-"`
}

func (o CreateTemplateGroupCollectionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTemplateGroupCollectionResponse struct{}"
	}

	return strings.Join([]string{"CreateTemplateGroupCollectionResponse", string(data)}, " ")
}
