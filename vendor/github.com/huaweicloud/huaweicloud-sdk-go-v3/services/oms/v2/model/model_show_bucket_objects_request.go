package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBucketObjectsRequest Request Object
type ShowBucketObjectsRequest struct {
	Body *ShowBucketReq `json:"body,omitempty"`
}

func (o ShowBucketObjectsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBucketObjectsRequest struct{}"
	}

	return strings.Join([]string{"ShowBucketObjectsRequest", string(data)}, " ")
}
