package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddActionRuleRequest Request Object
type AddActionRuleRequest struct {
	Body *ActionRule `json:"body,omitempty"`
}

func (o AddActionRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddActionRuleRequest struct{}"
	}

	return strings.Join([]string{"AddActionRuleRequest", string(data)}, " ")
}
