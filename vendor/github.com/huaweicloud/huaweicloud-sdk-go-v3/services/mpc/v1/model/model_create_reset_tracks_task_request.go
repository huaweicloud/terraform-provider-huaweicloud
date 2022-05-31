package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateResetTracksTaskRequest struct {
	Body *CreateResetTracksReq `json:"body,omitempty"`
}

func (o CreateResetTracksTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateResetTracksTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateResetTracksTaskRequest", string(data)}, " ")
}
