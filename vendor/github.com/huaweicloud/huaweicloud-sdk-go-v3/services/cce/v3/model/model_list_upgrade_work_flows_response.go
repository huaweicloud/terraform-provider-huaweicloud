package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUpgradeWorkFlowsResponse Response Object
type ListUpgradeWorkFlowsResponse struct {

	// API类型，固定值“List”，该值不可修改。
	Kind *string `json:"kind,omitempty"`

	// API版本，固定值“v3”，该值不可修改。
	ApiVersion *string `json:"apiVersion,omitempty"`

	Items          *UpgradeWorkFlow `json:"items,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListUpgradeWorkFlowsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUpgradeWorkFlowsResponse struct{}"
	}

	return strings.Join([]string{"ListUpgradeWorkFlowsResponse", string(data)}, " ")
}
