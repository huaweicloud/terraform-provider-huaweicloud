package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateMuteRuleResponse Response Object
type UpdateMuteRuleResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateMuteRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateMuteRuleResponse struct{}"
	}

	return strings.Join([]string{"UpdateMuteRuleResponse", string(data)}, " ")
}
