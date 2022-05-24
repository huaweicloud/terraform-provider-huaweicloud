package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateAssetCategoryRequest struct {
	Body *UpdateCategoryReq `json:"body,omitempty"`
}

func (o UpdateAssetCategoryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAssetCategoryRequest struct{}"
	}

	return strings.Join([]string{"UpdateAssetCategoryRequest", string(data)}, " ")
}
