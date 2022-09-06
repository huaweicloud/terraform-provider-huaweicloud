package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteTemplateGroupCollectionRequest struct {

	// 模板组集合id
	GroupCollectionId string `json:"group_collection_id"`
}

func (o DeleteTemplateGroupCollectionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTemplateGroupCollectionRequest struct{}"
	}

	return strings.Join([]string{"DeleteTemplateGroupCollectionRequest", string(data)}, " ")
}
