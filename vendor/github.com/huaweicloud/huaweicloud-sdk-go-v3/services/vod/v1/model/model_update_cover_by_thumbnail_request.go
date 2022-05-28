package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateCoverByThumbnailRequest struct {
	Body *UpdateCoverByThumbnailReq `json:"body,omitempty"`
}

func (o UpdateCoverByThumbnailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCoverByThumbnailRequest struct{}"
	}

	return strings.Join([]string{"UpdateCoverByThumbnailRequest", string(data)}, " ")
}
