package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateAssetMetaRequest struct {
	Body *UpdateAssetMetaReq `json:"body,omitempty"`
}

func (o UpdateAssetMetaRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAssetMetaRequest struct{}"
	}

	return strings.Join([]string{"UpdateAssetMetaRequest", string(data)}, " ")
}
