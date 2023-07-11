package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateObsBucketAuthorityPublicRequest Request Object
type UpdateObsBucketAuthorityPublicRequest struct {
	Body *ObsAuthorityConfigV2 `json:"body,omitempty"`
}

func (o UpdateObsBucketAuthorityPublicRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateObsBucketAuthorityPublicRequest struct{}"
	}

	return strings.Join([]string{"UpdateObsBucketAuthorityPublicRequest", string(data)}, " ")
}
