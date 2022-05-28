package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateBucketAuthorizedResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateBucketAuthorizedResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBucketAuthorizedResponse struct{}"
	}

	return strings.Join([]string{"UpdateBucketAuthorizedResponse", string(data)}, " ")
}
