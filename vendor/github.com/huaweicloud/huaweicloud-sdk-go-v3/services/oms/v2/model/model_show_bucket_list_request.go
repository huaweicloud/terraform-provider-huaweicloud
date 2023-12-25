package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBucketListRequest Request Object
type ShowBucketListRequest struct {
	Body *ListBucketsReq `json:"body,omitempty"`
}

func (o ShowBucketListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBucketListRequest struct{}"
	}

	return strings.Join([]string{"ShowBucketListRequest", string(data)}, " ")
}
