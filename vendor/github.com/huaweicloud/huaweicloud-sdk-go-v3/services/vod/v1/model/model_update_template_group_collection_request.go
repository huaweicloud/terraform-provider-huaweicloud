package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateTemplateGroupCollectionRequest struct {
	Body *ModifyTemplateGroupCollection `json:"body,omitempty"`
}

func (o UpdateTemplateGroupCollectionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTemplateGroupCollectionRequest struct{}"
	}

	return strings.Join([]string{"UpdateTemplateGroupCollectionRequest", string(data)}, " ")
}
