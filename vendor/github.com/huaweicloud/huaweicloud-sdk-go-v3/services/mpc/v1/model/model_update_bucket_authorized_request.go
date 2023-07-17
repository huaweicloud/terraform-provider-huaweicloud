package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateBucketAuthorizedRequest Request Object
type UpdateBucketAuthorizedRequest struct {
	Body *BucketAuthorizedReq `json:"body,omitempty"`
}

func (o UpdateBucketAuthorizedRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBucketAuthorizedRequest struct{}"
	}

	return strings.Join([]string{"UpdateBucketAuthorizedRequest", string(data)}, " ")
}
