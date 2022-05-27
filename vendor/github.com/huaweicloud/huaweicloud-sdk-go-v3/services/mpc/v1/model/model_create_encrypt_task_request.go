package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateEncryptTaskRequest struct {
	Body *CreateEncryptReq `json:"body,omitempty"`
}

func (o CreateEncryptTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateEncryptTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateEncryptTaskRequest", string(data)}, " ")
}
