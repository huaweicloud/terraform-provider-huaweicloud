package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteMuteRulesRequest struct {
	Body *[]DeleteMuteRuleName `json:"body,omitempty"`
}

func (o DeleteMuteRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMuteRulesRequest struct{}"
	}

	return strings.Join([]string{"DeleteMuteRulesRequest", string(data)}, " ")
}
