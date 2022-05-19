package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateThumbnailsTaskRequest struct {
	Body *CreateThumbReq `json:"body,omitempty"`
}

func (o CreateThumbnailsTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateThumbnailsTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateThumbnailsTaskRequest", string(data)}, " ")
}
