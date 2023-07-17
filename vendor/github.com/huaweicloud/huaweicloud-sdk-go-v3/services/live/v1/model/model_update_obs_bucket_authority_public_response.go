package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateObsBucketAuthorityPublicResponse Response Object
type UpdateObsBucketAuthorityPublicResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateObsBucketAuthorityPublicResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateObsBucketAuthorityPublicResponse struct{}"
	}

	return strings.Join([]string{"UpdateObsBucketAuthorityPublicResponse", string(data)}, " ")
}
