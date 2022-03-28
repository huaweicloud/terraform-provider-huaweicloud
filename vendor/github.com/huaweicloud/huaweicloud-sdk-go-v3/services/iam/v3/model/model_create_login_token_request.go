package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateLoginTokenRequest struct {
	Body *CreateLoginTokenRequestBody `json:"body,omitempty"`
}

func (o CreateLoginTokenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateLoginTokenRequest struct{}"
	}

	return strings.Join([]string{"CreateLoginTokenRequest", string(data)}, " ")
}
