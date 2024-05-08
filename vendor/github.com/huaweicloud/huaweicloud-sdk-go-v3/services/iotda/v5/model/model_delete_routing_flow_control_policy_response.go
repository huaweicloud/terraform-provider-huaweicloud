package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteRoutingFlowControlPolicyResponse Response Object
type DeleteRoutingFlowControlPolicyResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteRoutingFlowControlPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRoutingFlowControlPolicyResponse struct{}"
	}

	return strings.Join([]string{"DeleteRoutingFlowControlPolicyResponse", string(data)}, " ")
}
