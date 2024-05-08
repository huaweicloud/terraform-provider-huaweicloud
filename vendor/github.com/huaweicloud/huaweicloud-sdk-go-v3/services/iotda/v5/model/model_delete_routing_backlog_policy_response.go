package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteRoutingBacklogPolicyResponse Response Object
type DeleteRoutingBacklogPolicyResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteRoutingBacklogPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRoutingBacklogPolicyResponse struct{}"
	}

	return strings.Join([]string{"DeleteRoutingBacklogPolicyResponse", string(data)}, " ")
}
