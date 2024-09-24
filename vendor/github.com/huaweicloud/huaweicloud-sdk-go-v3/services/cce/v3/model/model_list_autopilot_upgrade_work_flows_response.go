package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotUpgradeWorkFlowsResponse Response Object
type ListAutopilotUpgradeWorkFlowsResponse struct {

	// API类型，固定值“List”，该值不可修改。
	Kind *string `json:"kind,omitempty"`

	// API版本，固定值“v3”，该值不可修改。
	ApiVersion *string `json:"apiVersion,omitempty"`

	Items          *UpgradeWorkFlow `json:"items,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListAutopilotUpgradeWorkFlowsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotUpgradeWorkFlowsResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotUpgradeWorkFlowsResponse", string(data)}, " ")
}
