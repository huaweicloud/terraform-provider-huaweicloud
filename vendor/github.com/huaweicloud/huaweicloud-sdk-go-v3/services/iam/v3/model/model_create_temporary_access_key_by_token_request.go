package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateTemporaryAccessKeyByTokenRequest struct {
	Body *CreateTemporaryAccessKeyByTokenRequestBody `json:"body,omitempty"`
}

func (o CreateTemporaryAccessKeyByTokenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTemporaryAccessKeyByTokenRequest struct{}"
	}

	return strings.Join([]string{"CreateTemporaryAccessKeyByTokenRequest", string(data)}, " ")
}
