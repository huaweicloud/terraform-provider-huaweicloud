package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateEventRuleResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateEventRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateEventRuleResponse struct{}"
	}

	return strings.Join([]string{"UpdateEventRuleResponse", string(data)}, " ")
}
