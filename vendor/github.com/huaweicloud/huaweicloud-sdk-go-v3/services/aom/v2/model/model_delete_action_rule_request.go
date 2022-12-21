package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteActionRuleRequest struct {
	Body *[]string `json:"body,omitempty"`
}

func (o DeleteActionRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteActionRuleRequest struct{}"
	}

	return strings.Join([]string{"DeleteActionRuleRequest", string(data)}, " ")
}
