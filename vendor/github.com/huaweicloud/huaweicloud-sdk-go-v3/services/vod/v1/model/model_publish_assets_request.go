package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type PublishAssetsRequest struct {
	Body *PublishAssetReq `json:"body,omitempty"`
}

func (o PublishAssetsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PublishAssetsRequest struct{}"
	}

	return strings.Join([]string{"PublishAssetsRequest", string(data)}, " ")
}
