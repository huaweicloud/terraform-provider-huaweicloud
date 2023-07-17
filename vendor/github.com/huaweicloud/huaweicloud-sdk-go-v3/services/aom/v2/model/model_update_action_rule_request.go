package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateActionRuleRequest Request Object
type UpdateActionRuleRequest struct {
	Body *ActionRule `json:"body,omitempty"`
}

func (o UpdateActionRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateActionRuleRequest struct{}"
	}

	return strings.Join([]string{"UpdateActionRuleRequest", string(data)}, " ")
}
