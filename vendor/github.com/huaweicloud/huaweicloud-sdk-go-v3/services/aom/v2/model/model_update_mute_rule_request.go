package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateMuteRuleRequest struct {
	Body *MuteRule `json:"body,omitempty"`
}

func (o UpdateMuteRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateMuteRuleRequest struct{}"
	}

	return strings.Join([]string{"UpdateMuteRuleRequest", string(data)}, " ")
}
