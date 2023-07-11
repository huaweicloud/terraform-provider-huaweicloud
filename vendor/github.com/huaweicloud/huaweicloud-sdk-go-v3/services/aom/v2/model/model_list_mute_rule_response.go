package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListMuteRuleResponse Response Object
type ListMuteRuleResponse struct {
	Body           *[]MuteRule `json:"body,omitempty"`
	HttpStatusCode int         `json:"-"`
}

func (o ListMuteRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMuteRuleResponse struct{}"
	}

	return strings.Join([]string{"ListMuteRuleResponse", string(data)}, " ")
}
