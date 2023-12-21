package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBucketRegionRequest Request Object
type ShowBucketRegionRequest struct {
	Body *ShowBucketRegionReq `json:"body,omitempty"`
}

func (o ShowBucketRegionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBucketRegionRequest struct{}"
	}

	return strings.Join([]string{"ShowBucketRegionRequest", string(data)}, " ")
}
