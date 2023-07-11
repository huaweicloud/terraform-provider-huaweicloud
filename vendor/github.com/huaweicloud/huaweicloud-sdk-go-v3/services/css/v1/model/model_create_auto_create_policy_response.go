package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAutoCreatePolicyResponse Response Object
type CreateAutoCreatePolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateAutoCreatePolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAutoCreatePolicyResponse struct{}"
	}

	return strings.Join([]string{"CreateAutoCreatePolicyResponse", string(data)}, " ")
}
